package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v2"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	// mqtt broker => mqtt_url := "mqtt://<user>:<pass>@1.227.57.115:1883/" + topic
	mqttUrl := config.MqttUrl + config.SubscribeTopic
	uri, err := url.Parse(mqttUrl)
	if err != nil {
		log.Fatal(err)
	}

	// kafka broker
	producer, err := newProducer(config.KafkaUrl)
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

	dataWindow := newDataWindow()

	client := connect(config.AssetId, uri)
	client.Subscribe(config.SubscribeTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		//log.Printf("[%s] comming: %s\n", config.AssetId, string(msg.Payload()))

		payload := newPayload()
		payload.UnmarshalJSON(msg.Payload())
		dataWindow.Add(*payload)

		//log.Printf("[%s] size: %d\n", config.AssetId, dataWindow.Size)
		if dataWindow.Size >= config.WindowSize {
			valueMap := make(map[string]interface{})

			eventTime := time.Now().UTC().Format(time.RFC3339) // 2019-01-12T01:02:03Z
			valueMap["event_time"] = eventTime
			valueMap["assetId"] = config.AssetId
			valueMap["sensorType"] = config.SensorType
			valueMap["sensorId"] = config.SensorId

			valueMap["velocity_rms"] = dataWindow.getRootMeanSquareOfVelocity()
			valueMap["distance_rms"] = dataWindow.getRootMeanSquareOfDistance()
			valueMap["velocity_mean"] = dataWindow.getMeanOfVelocity()
			valueMap["distance_mean"] = dataWindow.getMeanOfDistance()
			valueMap["velocity_min"] = dataWindow.getMinOfVelocity()
			valueMap["distance_min"] = dataWindow.getMinOfDistance()
			valueMap["velocity_max"] = dataWindow.getMaxOfVelocity()
			valueMap["distance_max"] = dataWindow.getMaxOfDistance()

			msgBytes, _ := json.Marshal(valueMap)
			message := string(msgBytes)

			msg := &sarama.ProducerMessage{
				Topic: config.ProduceTopic,
				Value: sarama.StringEncoder(message),
			}

			log.Printf("[%s] msg: %s\n", config.AssetId, message)

			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("[%s] FAILED to send message: %s\n", config.AssetId, err)
			} else {
				log.Printf("[%s] message sent to partition %d at offset %d\n", config.AssetId, partition, offset)
			}

			dataWindow.Reset()
		}
	})

	<-c
}

func connect(clientId string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
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
