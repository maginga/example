package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-co-op/gocron"
	"gopkg.in/yaml.v2"
)

func main() {
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

	s := gocron.NewScheduler(time.UTC)
	s.Every(config.ScheduleIntervalSec).Seconds().Do(func() {
		// data can be acquired
		sensorData := Call(config.RestUrl)

		for _, row := range sensorData.Data {
			valueMap := make(map[string]interface{})
			valueMap["assetId"] = config.AssetName
			valueMap["sensorType"] = config.SensorType
			valueMap["sensorId"] = config.SensorId

			timezone, _ := time.LoadLocation(config.LocalTimeZone)
			t, err := ParseIn(row.Time, timezone)
			if err != nil {
				log.Printf("err: %v", err)
			}
			valueMap["event_time"] = t.UTC().Format(time.RFC3339)

			for i, val := range row.Ai {
				paramName := "analog_input_" + strconv.Itoa(i)
				valueMap[paramName] = val
			}

			msgBytes, _ := json.Marshal(valueMap)
			message := string(msgBytes)
			msg := &sarama.ProducerMessage{
				Topic: config.Topic,
				Value: sarama.StringEncoder(message),
			}

			log.Printf("[%s] msg: %s\n", config.AssetName, message)

			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("%s > FAILED to send message: %s\n", config.AssetName, err)
			} else {
				log.Printf("%s > message sent to partition %d at offset %d\n", config.AssetName, partition, offset)
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
