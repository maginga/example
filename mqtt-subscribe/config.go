package main

type Config struct {
	MqttUrl string `yaml:"mqttUrl"`
	Topic   string `yaml:"topic"`
}
