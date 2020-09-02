package csv2json

import (
	"encoding/csv"
	"io"
	"os"
	"strings"

	"github.com/project-flogo/core/activity"
)

type Activity struct {
}

func init() {
	_ = activity.Register(&Activity{})
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Eval(context activity.Context) (done bool, err error) {
	logger := context.Logger()
	logger.Info("Executing CSV2JSON activity")

	input := &Input{}
	err = context.GetInputObject(input)
	if err != nil {
		return false, err
	}
	fileName := input.FileName
	csvfile, err := os.Open(fileName)
	if err != nil {

	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1

	rows := []string{}
	header := make([]string, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {

		}

		if len(header) == 0 {
			header = record
			// logger.Debug("csv headers: ", header)
		} else {
			var json string
			values := []string{}
			for i := range header {
				values = append(values, header[i]+":"+"\""+record[i]+"\"")
			}
			json = "{" + strings.Join(values, ",") + "}"
			rows = append(rows, json)
		}
	}

	logger.Info("rows: ", len(rows))
	output := &Output{Message: rows}
	err = context.SetOutputObject(output)
	if err != nil {
		return false, err
	}

	logger.Info("CSV2JSON activity completed")
	return true, nil
}
