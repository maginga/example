package filewatcher2

import (
	"github.com/project-flogo/core/data/coerce"
)

type HandlerSettings struct {
	DirName    string `md:"dirName"` // directory name for watching
	AssetName  string `md:"assetName"`
	SensorName string `md:"sensorName"`
	SensorType string `md:"sensorType"`
}

type Output struct {
	FileName   string `md:"fileName"` // file name to be changed.
	AssetName  string `md:"assetName"`
	SensorName string `md:"sensorName"`
	SensorType string `md:"sensorType"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"fileName":   o.FileName,
		"assetName":  o.AssetName,
		"sensorName": o.SensorName,
		"sensorType": o.SensorType,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.FileName, err = coerce.ToString(values["fileName"])
	o.AssetName, err = coerce.ToString(values["assetName"])
	o.SensorName, err = coerce.ToString(values["sensorName"])
	o.SensorType, err = coerce.ToString(values["sensorType"])

	if err != nil {
		return err
	}

	return nil
}
