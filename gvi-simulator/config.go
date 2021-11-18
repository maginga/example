package main

type Config struct {
	BrokerList          []string `yaml:"brokerList,flow"`
	Topic               string   `yaml:"topic"`
	ParamSpecTopic      string   `yaml:"paramSpecTopic"`
	ModelSpecTopic      string   `yaml:"modelSpecTopic"`
	ModelAlarmSpecTopic string   `yaml:"modelAlarmSpecTopic"`
	ParamAlarmSpecTopic string   `yaml:"paramAlarmSpecTopic"`
	FeatureSpecTopic    string   `yaml:"featureSpecTopic"`
	EventSpecTopic      string   `yaml:"eventSpecTopic"`
	MaxRetry            int      `yaml:"maxRetry"`
	IntervalMs          int      `yaml:"intervalMs"`
	AssetPrefix         string   `yaml:"assetPrefix"`
	AssetNumber         []string `yaml:"assetNumber"`
	FileName            string   `yaml:"fileName"`
	ParamSpecFile       string   `yaml:"paramSpecFile"`
	ModelSpecFile       string   `yaml:"modelSpecFile"`
	ModelAlarmSpecFile  string   `yaml:"modelAlarmSpecFile"`
	ParamAlarmSpecFile  string   `yaml:"paramAlarmSpecFile"`
	EventSpecFile       string   `yaml:"eventSpecFile"`
	FeatureSpecFile     string   `yaml:"featureSpecFile"`
}
