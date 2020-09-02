package csv2json

import (
	"github.com/project-flogo/core/data/coerce"
)

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

// type Output struct {
// 	JsonObject map[string]interface{} `md:"jsonObject"` // The data that sent kafka
// }

// func (o *Output) ToMap() map[string]interface{} {
// 	return map[string]interface{}{
// 		"jsonObject": o.JsonObject,
// 	}
// }

// func (o *Output) FromMap(values map[string]interface{}) error {

// 	var err error
// 	o.JsonObject, err = coerce.ToObject(values["jsonObject"])
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
