package main

type Config struct {
	MqttUrl        string   `yaml:"mqttUrl"`
	SubscribeTopic string   `yaml:"subscribeTopic"`
	KafkaUrl       []string `yaml:"kafkaUrl"`
	ProduceTopic   string   `yaml:"produceTopic"`
	AssetId        string   `yaml:"assetId"`
	SensorType     string   `yaml:"sensorType"`
	SensorId       string   `yaml:"sensorId"`
	WindowSize     int16    `yaml:"windowSize"`
}
