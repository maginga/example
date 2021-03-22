package main

import (
	"database/sql"
	"encoding/json"
	"example/rdbms-connector/rdb"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/Shopify/sarama"
	"github.com/carlescere/scheduler"
	"gopkg.in/yaml.v2"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

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

			go func(keyName string) {
				r1, err := db.Exec("UPDATE dbo.history SET sync02=? WHEREb sync02=? AND sn=?", 1, 0, keyName)
				if err != nil {
					log.Fatal(err)
				}
				n, err := r1.RowsAffected()
				if err != nil {
					panic(err)
				}
				log.Printf("[%s] %s rows were selected.\n", n, keyName)

				sqlStatement, err := rdb.NewSQLStatement(dbHelper, config.Query)
				if err != nil {
					panic(err)
				}

				params := make(map[string]interface{})
				params["status"] = 1
				params["assetName"] = keyName
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
				log.Printf("[%s] %s rows were returned.\n", n, keyName)

				for _, rowMap := range rowList {

					if t, ok := rowMap["ts"].(time.Time); ok {
						rowMap["event_time"] = t.UTC().Format(time.RFC3339)
					} else {
						//rowMap["event_time"] = time.Now().UTC().Format(time.RFC3339) // 2019-01-12T01:02:03Z
					}

					rowMap["assetName"] = keyName
					rowMap["sensorId"] = config.SensorName
					rowMap["sensorName"] = config.SensorName
					rowMap["sensorType"] = config.SensorType

					mapString, _ := json.Marshal(rowMap)
					message := string(mapString)

					msg := &sarama.ProducerMessage{
						Topic: config.Topic,
						Value: sarama.StringEncoder(message),
					}

					partition, offset, err := producer.SendMessage(msg)
					if err != nil {
						log.Printf("[%s] FAILED to send message: %s\n", keyName, err)
					} else {
						log.Printf("[%s] message: %s\n", keyName, message)
						log.Printf("[%s] message sent to partition %d at offset %d\n", keyName, partition, offset)
					}
					time.Sleep(time.Millisecond * 50)
				}

				r2, err := db.Exec("UPDATE dbo.history SET sync02=? WHERE sync02=? AND sn=?", 2, 1, keyName)
				if err != nil {
					log.Fatal(err)
				}
				rowaffected, err := r2.RowsAffected()
				if err != nil {
					panic(err)
				}
				log.Printf("[%s] sent %s rows.\n", rowaffected)

			}(assetName)
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
	runtime.Goexit()
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
			default:
				log.Printf("type unknown\n")
			}
		}

		//todo do we need to do column mapping
		rowList = append(rowList, resMap)
	}

	return rowList, rows.Err()
}
