package main

type Config struct {
	BrokerUrl []string `yaml:"brokerUrl"`
	AssetList []string `yaml:"assetList"`
	Topic     string   `yaml:"topic"`
	FilePath  string   `yaml:"filePath"`
	Loop      int      `yaml:"loop"`
}
