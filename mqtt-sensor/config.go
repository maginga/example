package main

type Config struct {
	MqttUrl  string `yaml:"mqttUrl"`
	Topic    string `yaml:"topic"`
	FilePath string `yaml:"filePath"`
	Loop     int    `yaml:"loop"`
	Serial   int    `yaml:"serial"`
}
