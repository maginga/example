package csvreader

import "github.com/project-flogo/core/data/coerce"

// Settings setting struct
type Settings struct {
	TimeColumnIndex   string `md:"timeColumnIndex"` // The time column in CSV file.
	ExcludeColumns    string `md:"excludeColumns"`  // Columns to be exclude.
	PhysicalAssetName string `md:"physicalAssetName"`
	SensorName        string `md:"sensorName"`
	SensorType        string `md:"sensorType"`
}

// Input input
type Input struct {
	FileName string `md:"fileName"` //
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"fileName": i.FileName,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	i.FileName, err = coerce.ToString(values["fileName"])
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
