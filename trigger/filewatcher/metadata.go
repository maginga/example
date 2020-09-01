package filewatcher

import (
	"github.com/project-flogo/core/data/coerce"
)

type HandlerSettings struct {
	DirName string `md:"dirName"` // directory name for watching
}

type Output struct {
	FileName string `md:"fileName"` // file name to be changed.
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"fileName": o.FileName,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.FileName, err = coerce.ToString(values["fileName"])
	if err != nil {
		return err
	}

	return nil
}
