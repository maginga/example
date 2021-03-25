package main

import (
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
	"syscall"
	"time"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/Shopify/sarama"
	"github.com/carlescere/scheduler"
	"gopkg.in/yaml.v2"
)

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
		// DB is safe to be used by multiple goroutines
		for _, assetName := range config.AssetList {

			sql1 := fmt.Sprintf("UPDATE dbo.history SET sync02=1 WHERE sync02=0 AND sn='%s' AND ts >= '%s'", assetName, config.StartTime)
			r1, err := db.Exec(sql1)
			if err != nil {
				log.Fatal(err)
			}
			n, err := r1.RowsAffected()
			if err != nil {
				panic(err)
			}
			if n <= 0 {
				return
			}
			log.Printf("[%s] %d rows were selected.\n", assetName, n)

			sqlStatement, err := rdb.NewSQLStatement(dbHelper, config.Query)
			if err != nil {
				panic(err)
			}

			params := make(map[string]interface{})
			params["status"] = 1
			params["assetName"] = assetName
			params["startTime"] = config.StartTime

			rows, err := db.Query(sqlStatement.ToStatementSQL(params))
			if err != nil {
				panic(err)
			}
			defer rows.Close()

			rowList, err := getLabeledResults(dbHelper, rows)
			if err != nil {
				panic(err)
			}
			log.Printf("[%s] %d rows ware returned.\n", assetName, len(rowList))

			go func(keyName string, rows []map[string]interface{}) {
				for _, rowMap := range rows {
					valueMap := make(map[string]interface{})
					for k, v := range rowMap {
						if k == "ts" || k == "tz" || k == "sn" || k == "sync01" || k == "sync02" || k == "sync03" || k == "id" || k == "ip" {
							continue
						}

						if v == nil {
							continue
						}

						valueMap[k] = v
					}

					if t1, ok := rowMap["tz"].(time.Time); ok {
						z, _ := t1.Zone()
						t2 := t1.UTC()

						if config.LogMessage {
							log.Printf("[%s] ZONE : %s, Local : %s, UTC: %s\n", keyName, z, t1, t2)
						}
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

					if config.LogMessage {
						if err != nil {
							log.Printf("[%s] FAILED to send message: %s\n", keyName, err)
						} else {
							log.Printf("[%s] message sent to partition %d at offset %d\n", keyName, partition, offset)
						}
					}

					time.Sleep(time.Duration(config.DelayMs) * time.Millisecond)
				}
			}(assetName, rowList)

			sql2 := fmt.Sprintf("UPDATE dbo.history SET sync02=2 WHERE sync02=1 AND sn='%s' AND ts >= '%s'", assetName, config.StartTime)
			r2, err := db.Exec(sql2)
			if err != nil {
				log.Fatal(err)
			}
			rowaffected, err := r2.RowsAffected()
			if err != nil {
				panic(err)
			}
			log.Printf("[%s] %d rows were sent.\n", assetName, rowaffected)

		}
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
			fmt.Printf("shutdown now.")
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
