package main

type Config struct {
	BrokerList []string `yaml:"brokerList,flow"`
	Topic      string   `yaml:"topic"`
	MaxRetry   int      `yaml:"maxRetry"`
	Assets     int      `yaml:"assets"`
	FileName   string   `yaml:"fileName"`
}
