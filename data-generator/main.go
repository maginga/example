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

	c := sarama.NewConfig()
	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Producer.Retry.Max = config.MaxRetry
	c.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(config.BrokerList, c)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panic(err)
		}
	}()

	for loop := 0; loop < 100000; loop++ {
		var wait sync.WaitGroup
		wait.Add(len(config.AssetNumber))

		msgList := mapTo(config.FileName)

		for _, assetNum := range config.AssetNumber {

			assetName := config.AssetPrefix + assetNum
			sensorId := "P" + assetNum

			go func(assetName, sensorId, topic string, p *sarama.SyncProducer, messages []string) {
				defer wait.Done()

				for _, msg := range messages {
					eventTime := time.Now().UTC().Format(time.RFC3339) // 2019-01-12T01:02:03Z
					val := strings.Replace(msg, "$1", eventTime, 1)
					val = strings.Replace(val, "$2", assetName, 1)
					val = strings.Replace(val, "$3", sensorId, 1)

					log.Printf("[%s] - %s\n", assetName, val)
					msg := &sarama.ProducerMessage{
						Topic: topic,
						Value: sarama.StringEncoder(val),
					}

					(*p).SendMessage(msg)
					//log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
					//log.Println("")

					time.Sleep(time.Second * 1)
				}

			}(assetName, sensorId, config.Topic, &producer, msgList)

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
		valueMap["assetId"] = "$2"
		valueMap["sensorId"] = "$3"
		valueMap["sensorType"] = "IFM"

		for j, val := range row {
			f, _ := strconv.ParseFloat(val, 64)
			valueMap[header[j]] = f
		}

		mapString, _ := json.Marshal(valueMap)
		jsonMsg := string(mapString)
		ret = append(ret, jsonMsg)
	}

	return ret
}
