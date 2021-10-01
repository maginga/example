package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-co-op/gocron"
	"gopkg.in/yaml.v2"
)

type SensorData struct {
	Data []DataValue
}

type DataValue struct {
	Ai    []float32
	Aiabh []bool
	Aiabl []bool
	Aialh []bool
	Aiall []bool
	Di    []bool
	Diab  []bool
	Dial  []bool
	Time  string
}

func main() {
	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	if config.Points < 1 {
		log.Print("The number of data points must be at least one.")
		os.Exit(0)
	}

	producer, err := newProducer(config.BrokerUrl)
	if err != nil {
		// Should not reach here
		panic(err)
	}

	lastTime := time.Time{}

	schedTime := config.Points - 1
	if config.Points <= 1 {
		schedTime = 1
	}

	s := gocron.NewScheduler(time.UTC)
	s.Every(schedTime).Seconds().Do(func() {
		// data can be acquired
		point := strconv.Itoa(config.Points)
		bodyBytes, err := GetBody(config.RestUrl + point)
		if err != nil {
			log.Println("unable to retrieve the response body from the Purewafer API server: %v", err)
		} else {
			var sensorData SensorData
			json.Unmarshal(bodyBytes, &sensorData)

			if sensorData.Data != nil {
				for i := range sensorData.Data {
					value := sensorData.Data[len(sensorData.Data)-1-i]

					timezone, _ := time.LoadLocation(config.LocalTimeZone)
					sensorTime, err := ParseIn(value.Time, timezone)
					if err != nil {
						log.Printf("err: %v", err)
					}

					if lastTime.Before(sensorTime) {
						valueMap := make(map[string]interface{})
						valueMap["assetId"] = config.AssetName
						valueMap["sensorType"] = config.SensorType
						valueMap["sensorId"] = config.SensorId
						valueMap["event_time"] = sensorTime.UTC().Format(time.RFC3339)

						for i, val := range value.Ai {
							paramName := "analog_input_" + strconv.Itoa(i+1)
							valueMap[paramName] = val
						}

						msgBytes, _ := json.Marshal(valueMap)
						message := string(msgBytes)
						msg := &sarama.ProducerMessage{
							Topic: config.Topic,
							Value: sarama.StringEncoder(message),
						}

						log.Printf("[%s] [%d] msg: %s\n", config.AssetName, i, message)

						partition, offset, err := producer.SendMessage(msg)
						if err != nil {
							log.Printf("[%s] > FAILED to send message: %s\n", config.AssetName, err)
						} else {
							log.Printf("[%s] > message sent to partition %d at offset %d\n", config.AssetName, partition, offset)
						}

						// for kafka
						time.Sleep(time.Millisecond * 200)
					} else {
						log.Print("This message has already been delivered.")
					}

					lastTime = sensorTime
				}
			} else {
				log.Print("There's no response.")
			}
		}
	})
	// starts the scheduler and blocks current execution path
	s.StartBlocking()

}

func newProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	// The level of acknowledgement reliability needed from the broker.
	producer, err := sarama.NewSyncProducer(brokers, config)
	// producer, err := sarama.NewAsyncProducer(brokers, config)

	return producer, err
}
