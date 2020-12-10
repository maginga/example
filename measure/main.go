package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
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

	c := sarama.NewConfig()
	c.Consumer.Return.Errors = true
	//c.Consumer.Group.Session.Timeout = 120
	c.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	brokers := config.BrokerList
	master, err := sarama.NewConsumer(brokers, c)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			log.Panic(err)
		}
	}()

	consumer, err := master.ConsumePartition(config.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Panic(err)
	}
	log.Println("comsuming...")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				// *messageCountStart++
				var mapdata map[string]interface{}
				if err := json.Unmarshal(msg.Value, &mapdata); err != nil {
					panic(err)
				}

				currentTime := time.Now().UTC()
				v1 := fmt.Sprintf("%.0f", mapdata["timestamp"]) //1.607500028e+12
				log.Println("timestamp: ", v1)

				v2, _ := strconv.ParseInt(v1, 10, 64)
				unixTimeUTC := time.Unix(0, v2*int64(time.Millisecond)) //gives unix time stamp in utc
				oldTime := unixTimeUTC.UTC()
				diff := currentTime.Sub(oldTime)

				log.Println("event     time: ", oldTime)
				log.Println("processed time: ", currentTime)

				log.Println("elapsed time: ", diff.Seconds(), " (sec)")

			case <-signals:
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	log.Println("Processed")
}
