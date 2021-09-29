package main

type Config struct {
	BrokerUrl []string `yaml:"brokerUrl"`
	Topic     string   `yaml:"topic"`
	LogFolder string   `yaml:"logFolder"`
	AssetId   string   `yaml:"assetId"`
}
