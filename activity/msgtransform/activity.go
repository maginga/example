package msgtransform

import (
	"encoding/json"
	"strconv"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})

// Activity define activity object
type Activity struct {
	settings *Settings
}

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
	ctx.Logger().Info("Executing transform activity")

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return false, err
	}

	message := input.Message

	var mapData map[string]interface{}
	if err := json.Unmarshal([]byte(message), &mapData); err != nil {
		ctx.Logger().Error(err)
	}

	// strict type
	var newMapData map[string]interface{}
	for key, value := range mapData {
		str, _ := coerce.ToString(value)
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			newMapData[key] = value
		} else {
			newMapData[key] = f
		}
	}

	newMapData["assetName"] = input.AssetName
	newMapData["sensorName"] = input.SensorName
	newMapData["sensorType"] = input.SensorType

	trans, _ := json.Marshal(newMapData)
	newMsg := string(trans)
	ctx.Logger().Infof("sending message: %v", newMsg)

	output := &Output{Message: newMsg}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}

	ctx.Logger().Info("transform activity completed")
	return true, nil
}
