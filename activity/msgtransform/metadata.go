package msgtransform

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings setting struct
type Settings struct {
	Option string `md:"option"` //
}

// Input input
type Input struct {
	Message    string `md:"message"` //
	AssetName  string `md:"assetName"`
	SensorName string `md:"sensorName"`
	SensorType string `md:"sensorType"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message":    i.Message,
		"assetName":  i.AssetName,
		"sensorName": i.SensorName,
		"sensorType": i.SensorType,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	i.Message, err = coerce.ToString(values["message"])
	i.AssetName, err = coerce.ToString(values["assetName"])
	i.SensorName, err = coerce.ToString(values["sensorName"])
	i.SensorType, err = coerce.ToString(values["sensorType"])

	if err != nil {
		return err
	}
	return nil
}

type Output struct {
	Message string `md:"message"` // Documents the partition that the message was placed on
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": o.Message,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Message, err = coerce.ToString(values["message"])
	if err != nil {
		return err
	}

	return nil
}
