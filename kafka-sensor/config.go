package main

type Config struct {
	BrokerUrl     []string `yaml:"brokerUrl"`
	AssetList     []string `yaml:"assetList"`
	Topic         string   `yaml:"topic"`
	FilePath      string   `yaml:"filePath"`
	SensorId      string   `yaml:"sensorId"`
	SensorType    string   `yaml:"sensorType"`
	Interval      string   `yaml:"interval"`
	Loop          int      `yaml:"loop"`
	StartIdx      int      `yaml:"startIdx"`
	Lot           int      `yaml:"lot"`
	Substrate     int      `yaml:"substrate"`
	SubstType     int      `yaml:"substType"`
	SubstLocation int      `yaml:"substLocation"`
	ProcessJob    int      `yaml:"processJob"`
	ControlJob    int      `yaml:"controlJob"`
	RecipeName    int      `yaml:"recipeName"`
}
