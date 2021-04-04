package main

import (
	"database/sql"
	"encoding/json"
	"example/pfeiffer-connector/rdb"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/Shopify/sarama"
	"github.com/carlescere/scheduler"
	"gopkg.in/yaml.v2"
)

var watermark sync.Map

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	watermark := sync.Map{}
	for _, name := range config.AssetList {
		watermark.Store(name, config.StartTime)
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
		log.Println("job started.")
		wg := sync.WaitGroup{}
		ch := make(chan string, len(config.AssetList))

		// DB is safe to be used by multiple goroutines
		for _, assetName := range config.AssetList {
			wg.Add(1)

			go func(assetName string) {
				defer wg.Done()

				fromTime, _ := watermark.Load(assetName)
				ft, _ := time.Parse("2006-01-02 15:04:05", fromTime.(string))
				tt := ft.Add(time.Duration(config.ChunkSizeMin) * time.Minute)
				toTime := tt.Format("2006-01-02 15:04:05")

				sqlStatement, err := rdb.NewSQLStatement(dbHelper, config.Query)
				if err != nil {
					log.Printf("[%s] %v\n", assetName, err)
					return
				}
				params := make(map[string]interface{})
				params["status"] = 0
				params["assetName"] = assetName
				params["fromTime"] = fromTime
				params["toTime"] = toTime

				selQuery := sqlStatement.ToStatementSQL(params)
				rows, err := db.Query(selQuery)
				if err != nil {
					log.Printf("[%s] %v\n", assetName, err)
					return
				}
				defer rows.Close()
				log.Printf("[%s] - rows are selected. (from: %s, to: %s)\n", assetName, fromTime, toTime)

				columns, err := rows.Columns()
				if err != nil {
					log.Printf("[%s] %v\n", assetName, err)
					return
				}

				columnTypes, err := rows.ColumnTypes()
				if err != nil {
					log.Panic(err)
				}

				var lastTime time.Time
				rowIndex := 0
				for rows.Next() {
					values := make([]interface{}, len(columnTypes))
					for i := range values {
						values[i] = dbHelper.GetScanType(columnTypes[i])
					}

					err = rows.Scan(values...)
					if err != nil {
						log.Panic(err)
					}

					valueMap := make(map[string]interface{})
					for i, column := range columns {
						//Exclude unnecessary columns.
						if column == "ts" || column == "sn" ||
							column == "sync01" || column == "sync02" || column == "sync03" ||
							column == "id" || column == "ip" {
							continue
						}
						switch v := values[i].(type) {
						case *interface{}:
							if column == "tz" {
								timeValue := *(v)
								lastTime = timeValue.(time.Time)
								utcTime := lastTime.UTC()
								valueMap["event_time"] = utcTime.Format(time.RFC3339)
							} else {
								if *(v) != nil {
									valueMap[column] = *(v)
								}
							}
						case *string:
							valueMap[column] = *(v)
						case *int:
							valueMap[column] = *(v)
						case *float32:
							valueMap[column] = *(v)
						default:
							log.Printf("type unknown\n")
						}
					}

					valueMap["assetId"] = assetName
					valueMap["sensorId"] = assetName
					valueMap["sensorName"] = assetName
					valueMap["sensorType"] = config.SensorType

					mapString, _ := json.Marshal(valueMap)
					message := string(mapString)

					if config.LogMessage {
						log.Printf("[%s] message: %s\n", assetName, message)
					}

					msg := &sarama.ProducerMessage{
						Topic: config.Topic,
						Value: sarama.StringEncoder(message),
					}

					partition, offset, err := producer.SendMessage(msg)
					rowIndex++

					if err != nil {
						log.Printf("[%s] FAILED to send message: %s\n", assetName, err)
					} else {
						log.Printf("[%s] sent (partition: %d,  offset: %d) - row: %d, (origin: %s, utc: %s)\n",
							assetName, partition, offset, rowIndex, lastTime.Format(time.RFC3339), valueMap["event_time"])
					}
					time.Sleep(time.Duration(config.DelayMs) * time.Millisecond)
				}

				if rowIndex <= 0 {
					return
				}

				resp := []string{assetName, fromTime.(string), toTime, lastTime.Format("2006-01-02 15:04:05")}
				ch <- strings.Join(resp, ",")

			}(assetName)
		}
		wg.Wait()
		close(ch)

		for chValue := range ch {
			values := strings.Split(chValue, ",")
			assetName := values[0]
			fromTime := values[1]
			toTime := values[2]
			lastTime := values[3]

			sql := fmt.Sprintf("UPDATE dbo.history SET sync02=2 WHERE sync02=0 AND sn='%s' AND ts >= '%s' AND ts < '%s'", assetName, fromTime, toTime)
			r, err := db.Exec(sql)
			if err != nil {
				log.Panic(err)
			}
			rowAffected, _ := r.RowsAffected()
			log.Printf("[%s] - [%d] rows are updated. (sync02: 2, from: %s, to: %s)\n", assetName, rowAffected, fromTime, toTime)

			watermark.Store(assetName, lastTime)
		}
		log.Println("job finished.")
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
