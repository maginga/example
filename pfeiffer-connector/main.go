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

var baseTime string

func loadConfig() (Config, error) {
	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	return config, err
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	config, err := loadConfig()
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

	baseTime := config.StartTime

	job := func() {
		log.Printf("job started. (base: %s, assets: %d)\n", baseTime, len(config.AssetList))

		sql1 := fmt.Sprintf("UPDATE dbo.history SET sync02=1 WHERE sync02=0 AND ts >= '%s' ", baseTime)
		sql1 += "AND ip IN ('" + strings.Join(config.IpAddress, "','") + "')"
		r1, err := db.Exec(sql1)
		if err != nil {
			log.Panic(err)
		}
		rowAffected1, _ := r1.RowsAffected()
		log.Printf("[%d] rows are changed. (state: 0 => 1)\n", rowAffected1)

		wg := sync.WaitGroup{}

		// DB is safe to be used by multiple goroutines
		for i, assetName := range config.AssetList {
			wg.Add(1)

			go func(idx int, assetName string) {
				defer wg.Done()

				sqlStatement, err := rdb.NewSQLStatement(dbHelper, config.Query)
				if err != nil {
					log.Panic(err)
				}
				params := make(map[string]interface{})
				params["status"] = 1
				params["ip"] = config.IpAddress[idx]
				params["startTime"] = baseTime

				selQuery := sqlStatement.ToStatementSQL(params)
				rows, err := db.Query(selQuery)
				if err != nil {
					log.Panic(err)
				}
				defer rows.Close()

				columns, err := rows.Columns()
				if err != nil {
					log.Panic(err)
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

					if config.Debug {
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
						if config.Debug {
							log.Printf("[%s][%s] sent (partition: %d,  offset: %d) - row: %d, (origin: %s, utc: %s)\n",
								assetName, config.IpAddress[idx], partition, offset, rowIndex, lastTime.Format(time.RFC3339), valueMap["event_time"])
						}
					}
					time.Sleep(time.Duration(config.DelayMs) * time.Millisecond)
				}

				log.Printf("[%s][%s] - [%d] rows sent. (last: %s)\n",
					assetName, config.IpAddress[idx], rowIndex, lastTime.Format(time.RFC3339))

				if rowIndex <= 0 {
					return
				}

			}(i, assetName)
		}
		wg.Wait()

		sql2 := fmt.Sprintf("UPDATE dbo.history SET sync02=2 WHERE sync02=1 AND ts >= '%s' ", baseTime)
		sql2 += "AND ip IN ('" + strings.Join(config.IpAddress, "','") + "')"
		r2, err := db.Exec(sql2)
		if err != nil {
			log.Panic(err)
		}
		rowAffected2, _ := r2.RowsAffected()
		log.Printf("[%d] rows are changed. (state: 1 => 2)\n", rowAffected2)

		// change the base time
		layout := "2006-01-02 15:04:05"
		utc, _ := time.Parse(layout, baseTime)
		loc, _ := time.LoadLocation(config.Timezone)
		localTime := utc.In(loc)
		days := time.Now().Sub(localTime).Hours() / 24
		if days > 7.0 {
			baseTime = localTime.AddDate(0, 0, 5).Format(layout)
			log.Printf("base time changed. (before: %s, after: %s)\n", localTime.Format(layout), baseTime)
		}

		log.Printf("job finished. (base diff: %f)\n", days)
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

	// Wait for SIGTERM
	waitForShutdown()
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

// waitForShutdown blocks until a SIGINT or SIGTERM is received.
func waitForShutdown() {
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	<-quit
}
