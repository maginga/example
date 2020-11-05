package csvsender

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})

// Activity define activity object
type Activity struct {
	conn     *KafkaConnection
	settings *Settings
}

// New create a new kafka activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	conn, err := getKafkaConnection(ctx.Logger(), s)
	if err != nil {
		ctx.Logger().Errorf("Kafka parameters initialization got error: [%s]", err.Error())
		return nil, err
	}

	act := &Activity{conn: conn, settings: s}
	return act, nil
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval evaluate
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	ctx.Logger().Info("Executing csvsender activity")

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

	pts, _ := time.ParseDuration(a.settings.PeriodOfTime)
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
				eventTime := time.Now().UTC().Format(time.RFC3339) // 2019-01-12T01:02:03Z
				valueMap["event_time"] = eventTime
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

			j, _ := json.Marshal(valueMap)
			jsonMsg := string(j)
			ctx.Logger().Infof("sending message: %v", jsonMsg)

			msg := &sarama.ProducerMessage{
				Topic: a.settings.Topic,
				Value: sarama.StringEncoder(jsonMsg),
			}
			partition, offset, err := a.conn.Connection().SendMessage(msg)
			if err != nil {
				return false, fmt.Errorf("failed to send Kakfa message for reason [%s]", err.Error())
			}

			ctx.Logger().Infof("A message sent on partition [%d] and offset [%d]", partition, offset)

			time.Sleep(pts)
		}
	}

	ctx.Logger().Info("csvsender activity completed")
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
