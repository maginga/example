package main

type Config struct {
	BrokerList []string `yaml:"brokerList,flow"`
	Topic      string   `yaml:"topic"`
	Partition  int      `yaml:"partition"`
	OffsetType int      `yaml:"offsetType"`
}
