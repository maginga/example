package filewatcher

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
}

type HandlerSettings struct {
	dirName string `md:"dirName"` // directory name for watching
}

type Output struct {
	filename string `md:"fileName"` // file name to be changed.
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"fileName": o.filename,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.filename, err = coerce.ToString(values["fileName"])
	if err != nil {
		return err
	}

	return nil
}
