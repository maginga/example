package main

type Config struct {
	BrokerList  []string `yaml:"brokerList,flow"`
	Topic       string   `yaml:"topic"`
	MaxRetry    int      `yaml:"maxRetry"`
	AssetPrefix string   `yaml:"assetPrefix"`
	AssetNumber []string `yaml:"assetNumber"`
	FileName    string   `yaml:"fileName"`
}
