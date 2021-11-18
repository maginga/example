package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
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

	producer, err := newProducer(config.BrokerList)
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

	// Open file for reading
	paramSpecFile, err := os.Open(config.ParamSpecFile)
	if err != nil {
		log.Fatal(err)
	}
	paramSpec, err := ioutil.ReadAll(paramSpecFile)
	if err != nil {
		log.Fatal(err)
	}

	// Open file for reading
	modelSpecFile, err := os.Open(config.ModelSpecFile)
	if err != nil {
		log.Fatal(err)
	}
	modelSpec, err := ioutil.ReadAll(modelSpecFile)
	if err != nil {
		log.Fatal(err)
	}

	// Open file for reading
	modelAlarmSpecFile, err := os.Open(config.ModelAlarmSpecFile)
	if err != nil {
		log.Fatal(err)
	}
	modelAlarmSpec, err := ioutil.ReadAll(modelAlarmSpecFile)
	if err != nil {
		log.Fatal(err)
	}

	// Open file for reading
	paramAlarmSpecFile, err := os.Open(config.ParamAlarmSpecFile)
	if err != nil {
		log.Fatal(err)
	}

	paramAlarmSpec, err := ioutil.ReadAll(paramAlarmSpecFile)
	if err != nil {
		log.Fatal(err)
	}

	featureSpecFile, err := os.Open(config.FeatureSpecFile)
	if err != nil {
		log.Fatal(err)
	}

	featureSpec, err := ioutil.ReadAll(featureSpecFile)
	if err != nil {
		log.Fatal(err)
	}

	eventSpecFile, err := os.Open(config.EventSpecFile)
	if err != nil {
		log.Fatal(err)
	}

	eventSpec, err := ioutil.ReadAll(eventSpecFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, assetNum := range config.AssetNumber {
		assetName := config.AssetPrefix + assetNum

		paramSpecvalue := strings.Replace(string(paramSpec), "$1", assetName, 1)
		msg1 := &sarama.ProducerMessage{
			Topic: config.ParamSpecTopic,
			Key:   sarama.StringEncoder(assetName),
			Value: sarama.StringEncoder(paramSpecvalue),
		}

		log.Printf("[%s] - %s\n", assetName, paramSpecvalue)

		_, offset, err := producer.SendMessage(msg1)
		if err != nil {
			log.Printf("%s > FAILED to send message: %s\n", assetName, err)
		} else {
			log.Printf("%s > sent parameter spec. at offset %d \n", assetName, offset)
		}

		modelSpecvalue := strings.Replace(string(modelSpec), "$1", assetName, 1)
		modelSpecvalue = strings.Replace(string(modelSpecvalue), "$2", assetNum, 1)
		msg2 := &sarama.ProducerMessage{
			Topic: config.ModelSpecTopic,
			Key:   sarama.StringEncoder(assetName),
			Value: sarama.StringEncoder(modelSpecvalue),
		}

		log.Printf("[%s] - %s\n", assetName, modelSpecvalue)

		_, offset, err = producer.SendMessage(msg2)
		if err != nil {
			log.Printf("%s > FAILED to send message: %s\n", assetName, err)
		} else {
			log.Printf("%s > sent model spec. at offset %d \n", assetName, offset)
		}

		modelAlarmSpecvalue := strings.Replace(string(modelAlarmSpec), "$1", assetName, 1)
		modelAlarmSpecvalue = strings.Replace(string(modelAlarmSpecvalue), "$2", assetNum, 1)
		msg3 := &sarama.ProducerMessage{
			Topic: config.ModelAlarmSpecTopic,
			Key:   sarama.StringEncoder(assetName),
			Value: sarama.StringEncoder(modelAlarmSpecvalue),
		}

		log.Printf("[%s] - %s\n", assetName, modelAlarmSpecvalue)

		_, offset, err = producer.SendMessage(msg3)
		if err != nil {
			log.Printf("%s > FAILED to send message: %s\n", assetName, err)
		} else {
			log.Printf("%s > sent model alarm spec. at offset %d \n", assetName, offset)
		}

		paramAlarmSpecvalue := strings.Replace(string(paramAlarmSpec), "$1", assetName, 1)
		msg4 := &sarama.ProducerMessage{
			Topic: config.ParamAlarmSpecTopic,
			Key:   sarama.StringEncoder(assetName),
			Value: sarama.StringEncoder(paramAlarmSpecvalue),
		}

		log.Printf("[%s] - %s\n", assetName, paramAlarmSpecvalue)

		_, offset, err = producer.SendMessage(msg4)
		if err != nil {
			log.Printf("%s > FAILED to send message: %s\n", assetName, err)
		} else {
			log.Printf("%s > sent parameter alarm spec. at offset %d \n", assetName, offset)
		}

		eventspecvalue := strings.Replace(string(eventSpec), "$1", assetName, 1)
		msg6 := &sarama.ProducerMessage{
			Topic: config.EventSpecTopic,
			Key:   sarama.StringEncoder(assetName),
			Value: sarama.StringEncoder(eventspecvalue),
		}

		log.Printf("[%s] - %s\n", assetName, eventspecvalue)

		_, offset, err = producer.SendMessage(msg6)
		if err != nil {
			log.Printf("%s > FAILED to send message: %s\n", assetName, err)
		} else {
			log.Printf("%s > sent parameter alarm spec. at offset %d \n", assetName, offset)
		}

		featurespecvalue := strings.Replace(string(featureSpec), "$1", assetName, 1)
		msg5 := &sarama.ProducerMessage{
			Topic: config.FeatureSpecTopic,
			Key:   sarama.StringEncoder(assetName),
			Value: sarama.StringEncoder(featurespecvalue),
		}

		log.Printf("[%s] - %s\n", assetName, featurespecvalue)

		_, offset, err = producer.SendMessage(msg5)
		if err != nil {
			log.Printf("%s > FAILED to send message: %s\n", assetName, err)
		} else {
			log.Printf("%s > sent parameter alarm spec. at offset %d \n", assetName, offset)
		}
	}

	for loop := 0; loop < 100000; loop++ {
		var wait sync.WaitGroup
		wait.Add(len(config.AssetNumber))

		msgList := mapTo(config.FileName)

		for _, assetNum := range config.AssetNumber {
			assetName := config.AssetPrefix + assetNum
			go func(assetName, topic string, p *sarama.SyncProducer, messages []string) {
				defer wait.Done()

				for _, msg := range messages {
					//eventTime := time.Now().UTC().Format(time.RFC3339) // 2006-01-02T15:04:05Z
					eventTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

					value := strings.Replace(msg, "$1", eventTime, 1)
					value = strings.Replace(value, "$2", assetName, 1)

					msg := &sarama.ProducerMessage{
						Topic: topic,
						Key:   sarama.StringEncoder(assetName),
						Value: sarama.StringEncoder(value),
					}

					log.Printf("[%s] - %s\n", assetName, value)

					partition, offset, err := producer.SendMessage(msg)
					if err != nil {
						log.Printf("%s > FAILED to send message: %s\n", assetName, err)
					} else {
						log.Printf("%s > message sent to partition %d at offset %d\n", assetName, partition, offset)
					}

					time.Sleep(time.Duration(config.IntervalMs) * time.Millisecond)
				}

			}(assetName, config.Topic, &producer, msgList)

		}
		wait.Wait()
	}
}

func mapTo(f string) []string {
	fileName, _ := filepath.Abs(f)
	csvfile, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1

	rows, _ := reader.ReadAll()

	header := make([]string, 0)
	ret := make([]string, 0)

	for i, row := range rows {
		valueMap := make(map[string]interface{})

		if i == 0 {
			if len(header) <= 0 {
				header = row
			}
			continue
		}

		valueMap["event_time"] = "$1"
		valueMap["asset_id"] = "$2"

		for j, val := range row {
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				valueMap[header[j]] = v
			} else {
				valueMap[header[j]] = val
			}
		}

		mapString, _ := json.Marshal(valueMap)
		jsonMsg := string(mapString)
		ret = append(ret, jsonMsg)
	}

	return ret
}

func newProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	// The level of acknowledgement reliability needed from the broker.
	producer, err := sarama.NewSyncProducer(brokers, config)
	// producer, err := sarama.NewAsyncProducer(brokers, config)

	return producer, err
}
