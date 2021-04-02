package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"example/rdbms-connector/rdb"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/Shopify/sarama"
	"github.com/carlescere/scheduler"
	"gopkg.in/yaml.v2"
)

var watermark string

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGUSR1)

	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	watermark := config.StartTime

	producer, err := newProducer(config.BrokerUrl)
	if err != nil {
		// Should not reach here
		panic(err)
	}
	log.Printf("Broker Address: %s\n", config.BrokerUrl)

	defer func() {
		if err := producer.Close(); err != nil {
			// Should not reach here
			panic(err)
		}
	}()

	log.Printf("MSSQL Conn: %s\n", config.ConnString)
	db, err := sql.Open("sqlserver", config.ConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbHelper, err := rdb.GetDbHelper()
	if err != nil {
		panic(err)
	}

	job := func() {
		var wg sync.WaitGroup

		// DB is safe to be used by multiple goroutines
		for _, assetName := range config.AssetList {
			// begin transaction
			ctx := context.Background()
			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				log.Fatal(err)
			}

			sql1 := fmt.Sprintf("UPDATE dbo.history SET sync02=1 WHERE sync02=0 AND sn='%s' AND ts >= '%s'", assetName, watermark)
			r1, err := tx.ExecContext(ctx, sql1)
			if err != nil {
				tx.Rollback()
				return
			}
			n, _ := r1.RowsAffected()

			if n <= 0 {
				err = tx.Commit()
				if err != nil {
					log.Println(err)
				}
				log.Printf("[%s] [%d] rows, skipped.\n", assetName, n)
				continue
			} else {
				log.Printf("[%s] [%d] rows, sql: %s\n", assetName, n, sql1)
			}

			wg.Add(1)

			sqlStatement, err := rdb.NewSQLStatement(dbHelper, config.Query)
			if err != nil {
				tx.Rollback()
				return
			}

			params := make(map[string]interface{})
			params["status"] = 1
			params["assetName"] = assetName
			params["startTime"] = watermark

			selQuery := sqlStatement.ToStatementSQL(params)
			rows, err := tx.Query(selQuery)
			if err != nil {
				tx.Rollback()
				return
			}
			defer rows.Close()

			rowList, err := getLabeledResults(dbHelper, rows)
			if err != nil {
				tx.Rollback()
				return
			}
			log.Printf("[%s] [%d] rows, sql: %s\n", assetName, len(rowList), selQuery)

			go func(keyName string, rows []map[string]interface{}) {
				defer wg.Done()

				totalRowCnt := len(rows)
				rowIndex := 0

				for _, rowMap := range rows {
					valueMap := make(map[string]interface{})
					for k, v := range rowMap {
						//Exclude unnecessary columns.
						if k == "ts" || k == "tz" || k == "sn" || k == "sync01" || k == "sync02" || k == "sync03" || k == "id" || k == "ip" {
							continue
						}

						//Exclude null values.
						if v == nil {
							continue
						}

						valueMap[k] = v
					}

					t1, ok := rowMap["tz"].(time.Time)
					if ok {
						t2 := t1.UTC()
						valueMap["event_time"] = t2.Format(time.RFC3339)
					} else {
						//valueMap["event_time"] = time.Now().UTC().Format(time.RFC3339) // 2019-01-12T01:02:03Z
					}

					valueMap["assetId"] = keyName
					valueMap["sensorId"] = keyName
					valueMap["sensorName"] = keyName
					valueMap["sensorType"] = config.SensorType

					mapString, _ := json.Marshal(valueMap)
					message := string(mapString)

					if config.LogMessage {
						log.Printf("[%s] message: %s\n", keyName, message)
					}

					msg := &sarama.ProducerMessage{
						Topic: config.Topic,
						Value: sarama.StringEncoder(message),
					}

					partition, offset, err := producer.SendMessage(msg)

					rowIndex++
					if err != nil {
						log.Printf("[%s] FAILED to send message: %s\n", keyName, err)
					} else {
						log.Printf("[%s] sent (partition: %d,  offset: %d) - (origin: %s, utc: %s) - (%d/%d)\n",
							keyName, partition, offset, t1.Format(time.RFC3339), valueMap["event_time"], rowIndex, totalRowCnt)
					}

					time.Sleep(time.Duration(config.DelayMs) * time.Millisecond)
				}
			}(assetName, rowList)

			sql2 := fmt.Sprintf("UPDATE dbo.history SET sync02=2 WHERE sync02=1 AND sn='%s' AND ts >= '%s'", assetName, watermark)
			r2, err := tx.ExecContext(ctx, sql2)
			if err != nil {
				tx.Rollback()
				return
			}
			rowaffected, _ := r2.RowsAffected()
			log.Printf("[%s] [%d] rows, sql: %s\n", assetName, rowaffected, sql2)

			// Commit the change if all queries ran successfully
			err = tx.Commit()
			if err != nil {
				log.Println(err)
			}
		}

		wg.Wait()
	}

	d, err := time.ParseDuration(config.RepeatInterval)
	if err != nil {
		log.Fatalf("unable to parse repeat interval: %s", err.Error())
		panic(err)
	}

	repeatInterval := int(d.Seconds())
	log.Printf("Scheduling action to repeat every %d seconds\n", repeatInterval)

	_, err = scheduler.Every(repeatInterval).Seconds().Run(job)
	if err != nil {
		log.Printf("Error scheduling repeating timer: ", err.Error())
	}

	// Keep the program from not exiting.
	//runtime.Goexit()

	for {
		signal := <-signalChan
		switch signal {
		case os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
			fmt.Printf("signal:%d\n", signal)
			os.Exit(1)
			break
		default:
			//fmt.Printf("Unknown signal(%d)\n", signal)
		}
	}
}

func newProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	// The level of acknowledgement reliability needed from the broker.
	return sarama.NewSyncProducer(brokers, config)
	// producer, err := sarama.NewAsyncProducer(brokers, config)
}

func prepareMessage(topic, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}

	return msg
}

func getLabeledResults(dbHelper rdb.DbHelper, rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var rowList []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columnTypes))
		for i := range values {
			values[i] = dbHelper.GetScanType(columnTypes[i])
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		resMap := make(map[string]interface{}, len(columns))
		for i, column := range columns {
			switch v := values[i].(type) {
			case *interface{}:
				resMap[column] = *(v)
			case *string:
				resMap[column] = *(v)
			case *int:
				resMap[column] = *(v)
			case *float32:
				resMap[column] = *(v)
			default:
				log.Printf("type unknown\n")
			}
		}

		//todo do we need to do column mapping
		rowList = append(rowList, resMap)
	}

	return rowList, rows.Err()
}
