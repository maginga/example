package csvsender

import "github.com/project-flogo/core/data/coerce"

// Settings setting struct
type Settings struct {
	TimeColumnIndex string `md:"timeColumnIndex"`     // The time column in CSV file.
	ExcludeColumns  string `md:"excludeColumns"`      // Columns to be exclude.
	BrokerUrls      string `md:"brokerUrls,required"` // The Kafka cluster to connect to
	User            string `md:"user"`                // If connecting to a SASL enabled port, the user id to use for authentication
	Password        string `md:"password"`            // If connecting to a SASL enabled port, the password to use for authentication
	TrustStore      string `md:"trustStore"`          // If connecting to a TLS secured port, the directory containing the certificates representing the trust chain for the connection. This is usually just the CACert used to sign the server's certificate
	Topic           string `md:"topic,required"`      // The Kafka topic on which to place the message
	PeriodOfTime    string `md:"periodOfTime"`        // ("s", "m", "h")
}

// Input input
type Input struct {
	FileName   string `md:"fileName"` //
	AssetName  string `md:"assetName"`
	SensorName string `md:"sensorName"`
	SensorType string `md:"sensorType"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"fileName":   i.FileName,
		"assetName":  i.AssetName,
		"sensorName": i.SensorName,
		"sensorType": i.SensorType,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	i.FileName, err = coerce.ToString(values["fileName"])
	i.AssetName, err = coerce.ToString(values["assetName"])
	i.SensorName, err = coerce.ToString(values["sensorName"])
	i.SensorType, err = coerce.ToString(values["sensorType"])

	if err != nil {
		return err
	}
	return nil
}

type Output struct {
	Partition int32 `md:"partition"` // Documents the partition that the message was placed on
	OffSet    int64 `md:"offset"`    // Documents the offset for the message
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"partition": o.Partition,
		"offset":    o.OffSet,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Partition, err = coerce.ToInt32(values["partition"])
	if err != nil {
		return err
	}

	o.OffSet, err = coerce.ToInt64(values["offset"])
	if err != nil {
		return err
	}

	return nil
}
