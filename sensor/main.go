package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v2"
)

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

func listen(uri *url.URL, topic string) {
	client := connect("sub", uri)
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
}

func main() {
	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	topic := config.Topic
	//mqtt_url := "mqtt://<user>:<pass>@1.227.57.115:1883/" + topic
	mqtt_url := config.MqttUrl + topic
	uri, err := url.Parse(mqtt_url)
	if err != nil {
		log.Fatal(err)
	}
	client := connect("pub", uri)

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
				valueMap[header[j]] = val
			}

			mapString, _ := json.Marshal(valueMap)
			jsonMsg := string(mapString)
			client.Publish(topic, 0, false, jsonMsg)
			log.Printf("event: %v", jsonMsg)

			time.Sleep(time.Second * 1)
		}
	}
}
