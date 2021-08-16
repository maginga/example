package main

type Config struct {
	LocalTimeZone string   `yaml:"localTimeZone"`
	RestUrl       string   `yaml:"restUrl"`
	Points        int      `yaml:"points"`
	BrokerUrl     []string `yaml:"brokerUrl"`
	Topic         string   `yaml:"topic"`
	AssetName     string   `yaml:"assetName"`
	SensorId      string   `yaml:"sensorId"`
	SensorType    string   `yaml:"sensorType"`
}
