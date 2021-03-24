package main

type Config struct {
	ConnString     string   `yaml:"connString"`
	RepeatInterval string   `yaml:"repeatInterval"`
	BrokerUrl      []string `yaml:"brokerUrl"`
	Topic          string   `yaml:"topic"`
	Query          string   `yaml:"query"`
	StartTime      string   `yaml:"startTime"`
	SensorType     string   `yaml:"sensorType"`
	AssetList      []string `yaml:"assetList"`
	DelayMs        int      `yaml:"delayMs"`
	LogMessage     bool     `yaml:"logMessage"`
	Location       string   `yaml:"location"`
}
