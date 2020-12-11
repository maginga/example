package csvreader

import "github.com/project-flogo/core/data/coerce"

// Settings setting struct
type Settings struct {
	TimeColumnIndex string `md:"timeColumnIndex"` // The time column in CSV file.
	TimeZone        string `md:"timeZone"`
	ExcludeColumns  string `md:"excludeColumns"` // Columns to be exclude.
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
	Message interface{} `md:"message"` // The data that sent kafka
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": o.Message,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Message, err = coerce.ToObject(values["message"])
	if err != nil {
		return err
	}

	return nil
}
