package binary2json

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

// Activity define activity object
type Activity struct {
	settings *Settings
}

type Header struct {
	totalSize uint32
	dataCount uint32
	datetime  string
	null      byte
}

type Record struct {
	data []float32
	cr   byte
	lf   byte
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
	logger.Info("Executing BINARY2JSON activity")

	input := &Input{}
	err = context.GetInputObject(input)
	if err != nil {
		return false, err
	}
	fileName := input.FileName
	file, err := os.Open(fileName)
	if err != nil {

	}
	defer file.Close()

	columns := make(map[string]interface{})
	for key, value := range a.settings.Columns {
		columns[key] = value
		logger.Info(key + ":" + fmt.Sprintf("%v", value))
	}

	headerSize := 32 // bytes
	// rows := []string{}
	for {
		header := Header{}
		headerBytes := readNextBytes(file, headerSize)

		buffer1 := bytes.NewBuffer(headerBytes)
		err = binary.Read(buffer1, binary.LittleEndian, &header)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		record := Record{}
		dataSize := int(header.totalSize) - headerSize
		dataBytes := readNextBytes(file, dataSize)

		buffer2 := bytes.NewBuffer(dataBytes)
		err = binary.Read(buffer2, binary.LittleEndian, &record)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		// var json string
		// values := []string{}

		for i := 1; i <= len(columns); i++ {
			logger.Info(fmt.Sprintf("%v", columns[strconv.Itoa(i)]) + ": " + fmt.Sprintf("%v", record.data[i]))
		}

		// for real := range record.data {
		// 	logger.Info(real)
		// }
		// json = "{" + strings.Join(values, ",") + "}"
		// rows = append(rows, json)
	}

	// logger.Info("rows: ", len(rows))
	// output := &Output{Message: rows}
	// err = context.SetOutputObject(output)
	// if err != nil {
	// 	return false, err
	// }

	logger.Info("BINARY2JSON activity completed")
	return true, nil
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
