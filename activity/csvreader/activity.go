package csvreader

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

// Activity define activity object
type Activity struct {
	settings *Settings
}

func init() {
	_ = activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})

// New create a new kafka activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	act := &Activity{settings: s}
	return act, nil
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval evaluate
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	ctx.Logger().Info("Executing csvreader activity")

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return false, err
	}
	fileName := input.FileName
	csvfile, err := os.Open(fileName)
	if err != nil {
		ctx.Logger().Errorf("err: %v", err)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1

	timeColIndex, err := strconv.Atoi(a.settings.TimeColumnIndex)
	if err != nil {
		timeColIndex = -1
	}

	var excludeColumns []int
	if len(a.settings.ExcludeColumns) > 0 {
		strs := strings.Split(a.settings.ExcludeColumns, ",")
		excludeColumns = make([]int, len(strs))
		for i := range excludeColumns {
			excludeColumns[i], _ = strconv.Atoi(strs[i])
		}
	}

	var rows []interface{}
	header := make([]string, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			ctx.Logger().Errorf("err: %v", err)
		}

		if len(header) == 0 {
			header = record
		} else {
			// var json string
			valueMap := make(map[string]interface{})

			if timeColIndex < 0 {
				valueMap["event_time"] = nil
			} else {
				t, _ := ParseLocal(record[timeColIndex])
				valueMap["event_time"] = t.Format(time.RFC3339)
			}

			valueMap["assetId"] = input.AssetName
			valueMap["sensorType"] = input.SensorType
			valueMap["sensorId"] = input.SensorName

			for i := range header {
				if excludeColumns != nil && contains(excludeColumns, i) {
					continue
				} else {
					if timeColIndex < 0 || timeColIndex != i {
						f1, e := strconv.ParseFloat(record[i], 8)
						if e == nil {
							valueMap[header[i]] = f1
						}
					}
				}
			}
			rows = append(rows, valueMap)
		}
	}

	ctx.Logger().Info("rows: ", len(rows))
	output := &Output{Message: rows}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}

	ctx.Logger().Info("csvreader activity completed")
	return true, nil
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
