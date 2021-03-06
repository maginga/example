package main

type Config struct {
	ConnString     string   `yaml:"connString"`
	RepeatInterval string   `yaml:"repeatInterval"`
	BrokerUrl      []string `yaml:"brokerUrl"`
	Topic          string   `yaml:"topic"`
	Query          string   `yaml:"query"`
	StartTime      string   `yaml:"startTime"`
	SensorType     string   `yaml:"sensorType"`
	IpAddress      []string `yaml:"ipAddress"`
	AssetList      []string `yaml:"assetList"`
	DelayMs        int      `yaml:"delayMs"`
	Debug          bool     `yaml:"debug"`
	Timezone       string   `yaml:"timezone"`
}
