package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
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

	defer func() {
		if err := producer.Close(); err != nil {
			// Should not reach here
			panic(err)
		}
	}()

	file, _ := os.Open(config.FilePath)
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.FieldsPerRecord = -1

	rows, _ := reader.ReadAll()

	header := make([]string, 0)
	valueMap := make(map[string]interface{})

	for l := 1; l <= config.Loop; l++ {

		for i, row := range rows {
			if i == 0 {
				if len(header) <= 0 {
					header = row
				}

				continue
			}

			eventTime := time.Now().UTC().Format(time.RFC3339) // 2019-01-12T01:02:03Z
			valueMap["event_time"] = eventTime

			for j, val := range row {
				valueMap[header[j]], _ = strconv.ParseFloat(val, 64)
			}

			for _, assetName := range config.AssetList {

				targetMap := make(map[string]interface{})
				for key, value := range valueMap {
					targetMap[key] = value
				}

				go func(asset string, valueMap map[string]interface{}) {

					valueMap["assetId"] = asset
					valueMap["sensorType"] = "Pump"
					valueMap["sensorName"] = asset
					valueMap["sensorId"] = asset

					mapString, _ := json.Marshal(valueMap)
					message := string(mapString)

					msg := &sarama.ProducerMessage{
						Topic: config.Topic,
						Value: sarama.StringEncoder(message),
					}

					log.Printf("[%s] msg: %s\n", asset, message)

					partition, offset, err := producer.SendMessage(msg)
					if err != nil {
						log.Printf("%s > FAILED to send message: %s\n", asset, err)
					} else {
						log.Printf("%s > message sent to partition %d at offset %d\n", asset, partition, offset)
					}
				}(assetName, targetMap)
			}

			time.Sleep(time.Second * 1)
		}
	}
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

func prepareMessage(topic, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}

	return msg
}
