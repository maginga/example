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

type Activity struct {
	timeColumnIndex int
	excludeColumns  []int
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

	timeIdx, err := strconv.Atoi(s.TimeColumnIndex)
	if err != nil {
		timeIdx = -1
	}

	var excludeColumns []int
	if len(s.ExcludeColumns) > 0 {
		strs := strings.Split(s.ExcludeColumns, ",")
		excludeColumns = make([]int, len(strs))
		for i := range excludeColumns {
			excludeColumns[i], _ = strconv.Atoi(strs[i])
		}
	}

	act := &Activity{timeColumnIndex: timeIdx, excludeColumns: excludeColumns}
	return act, nil
}

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
		} else {
			var json string
			values := []string{}

			if a.timeColumnIndex == -1 {
				timeStr := time.Now().UTC().Format(time.RFC3339) // 2019-01-12T01:02:03Z
				values = append(values, "event_time:\""+timeStr+"\"")
			} else {
				t, _ := ParseLocal(record[a.timeColumnIndex])
				values = append(values, "event_time:\""+t.Format(time.RFC3339)+"\"")
			}

			for i := range header {
				if contains(a.excludeColumns, i) {
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
