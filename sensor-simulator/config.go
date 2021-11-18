package main

type Config struct {
	BrokerUrl   []string `yaml:"brokerUrl"`
	AssetList   []string `yaml:"assetList"`
	SensorId    string   `yaml:"sensorId"`
	SensorType  string   `yaml:"sensorType"`
	Topic       string   `yaml:"topic"`
	FilePath    string   `yaml:"filePath"`
	Interval    string   `yaml:"interval"`
	Loop        int      `yaml:"loop"`
	StartColumn int      `yaml:"startColumn"`
}
