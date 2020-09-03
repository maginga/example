package csv2json

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
		} else {
			var json string
			values := []string{}

			if timeColIndex == -1 {
				timeStr := time.Now().UTC().Format(time.RFC3339) // 2019-01-12T01:02:03Z
				values = append(values, "event_time:\""+timeStr+"\"")
			} else {
				//t, _ := ParseLocal(record[timeColIndex])
				//values = append(values, "event_time:\""+t.Format(time.RFC3339)+"\"")
			}

			for i := range header {
				if excludeColumns != nil && contains(excludeColumns, i) {
					continue
				} else {
					values = append(values, header[i]+":"+"\""+record[i]+"\"")
				}
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

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
