package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

/*
  json file loader for inx site
*/
func main() {
	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	if len(config.LogFolder) <= 0 {
		log.Print("exit.")
		os.Exit(0)
	}

	if len(config.BrokerUrl) <= 0 {
		log.Print("exit.")
		os.Exit(0)
	}

	producer, err := newProducer(config.BrokerUrl)
	if err != nil {
		panic(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(config.LogFolder)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Watching : ", config.LogFolder)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					file, err := os.Open(event.Name)
					if err != nil {
						log.Fatalf("failed to open")
					}
					defer file.Close()

					scanner := bufio.NewScanner(file)
					scanner.Split(bufio.ScanLines)

					for scanner.Scan() {
						line := scanner.Text()

						msg := &sarama.ProducerMessage{
							Topic: config.Topic,
							Key:   sarama.StringEncoder(config.AssetId),
							Value: sarama.StringEncoder(line),
						}
						log.Printf("[%s] msg: %s\n", config.AssetId, line)

						partition, offset, err := producer.SendMessage(msg)
						if err != nil {
							log.Printf("[%s] > FAILED to send message: %s\n", config.AssetId, err)
						} else {
							log.Printf("[%s] > message sent to partition %d at offset %d\n", config.AssetId, partition, offset)
						}

						// for kafka
						time.Sleep(time.Millisecond * 200)
					}

					if scanner.Err() != nil {
						log.Println(scanner.Err())
					}
				}
			case err := <-watcher.Errors:
				log.Fatal("Error: ", err.Error())
			}
		}
	}()

	// Wait for SIGTERM
	waitForShutdown()
}

// waitForShutdown blocks until a SIGINT or SIGTERM is received.
func waitForShutdown() {
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	<-quit
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
