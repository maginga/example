package producer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})

// Activity is a kafka activity
type Activity struct {
	conn     *KafkaConnection
	settings *Settings
}

// New create a new kafka activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	conn, err := getKafkaConnection(ctx.Logger(), settings)
	if err != nil {
		ctx.Logger().Errorf("Kafka parameters initialization got error: [%s]", err.Error())
		return nil, err
	}

	act := &Activity{conn: conn, settings: settings}
	return act, nil
}

// Metadata returns the metadata for the kafka activity
func (*Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements the evaluation of the kafka activity
func (act *Activity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}

	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	if input.Message == "" {
		return false, fmt.Errorf("no message to publish")
	}

	pts, _ := time.ParseDuration(act.settings.PeriodOfTime)
	time.Sleep(pts)

	ctx.Logger().Infof("sending message: %v", input.Message)
	message := makeMessage(input.Message.(string))

	msg := &sarama.ProducerMessage{
		Topic: act.settings.Topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := act.conn.Connection().SendMessage(msg)
	if err != nil {
		return false, fmt.Errorf("failed to send Kakfa message for reason [%s]", err.Error())
	}

	output := &Output{}
	output.Partition = partition
	output.OffSet = offset

	// if ctx.Logger().DebugEnabled() {
	// 	ctx.Logger().Debugf("Kafka message [%v] sent successfully on partition [%d] and offset [%d]",
	// 		input.Message, partition, offset)
	// }

	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}

	ctx.Logger().Infof("A message sent on partition [%d] and offset [%d]", partition, offset)

	return true, nil
}

func makeMessage(record string) string {
	b := []byte(record)
	var m map[string]interface{}
	json.Unmarshal(b, &m)

	if m["event_time"] == nil {
		eventTime := time.Now().UTC().Format(time.RFC3339) // 2019-01-12T01:02:03Z
		m["event_time"] = eventTime
	}

	j, _ := json.Marshal(m)
	s := string(j)
	return s
}
