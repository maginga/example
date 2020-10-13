package binreader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
	logger.Info("Executing binreader activity")

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

	rows := []string{}
	for {
		var totalSize uint32
		tsBytes, e := readNextBytes(file, 4)
		if e != nil && e == io.EOF {
			break
		}

		tsBuffer := bytes.NewBuffer(tsBytes)
		err = binary.Read(tsBuffer, binary.LittleEndian, &totalSize)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		var dataCount uint32
		dcBytes, _ := readNextBytes(file, 4)
		dcBuffer := bytes.NewBuffer(dcBytes)
		err = binary.Read(dcBuffer, binary.LittleEndian, &dataCount)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		var dateTime [23]byte
		dtBytes, _ := readNextBytes(file, 23)
		dtBuffer := bytes.NewBuffer(dtBytes)
		err = binary.Read(dtBuffer, binary.LittleEndian, &dateTime)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}
		strDateTime := string(dateTime[:23]) //BytesToString(datetime)

		var null byte
		nullBytes, _ := readNextBytes(file, 1)
		nullBuffer := bytes.NewBuffer(nullBytes)
		err = binary.Read(nullBuffer, binary.LittleEndian, &null)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		values := []string{}
		itemCount := int(dataCount)
		for i := 0; i < itemCount; i++ {
			var dataPoint float32
			dpBytes, _ := readNextBytes(file, 4)
			dpBuffer := bytes.NewBuffer(dpBytes)
			err = binary.Read(dpBuffer, binary.LittleEndian, &dataPoint)
			if err != nil {
				log.Fatal("binary.Read failed", err)
			}
			// values = append(values, fmt.Sprintf("%v", columns[strconv.Itoa(i)])+"="+fmt.Sprintf("%f", dataPoint))
			values = append(values, fmt.Sprintf("%f", dataPoint))
		}
		json := "{event_time=" + strDateTime + "," + strings.Join(values, ",") + "}"
		rows = append(rows, json)

		var crlf []byte
		crlfBytes, _ := readNextBytes(file, 2)
		crlfBuffer := bytes.NewBuffer(crlfBytes)
		err = binary.Read(crlfBuffer, binary.LittleEndian, &crlf)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}
	}

	logger.Info("rows: ", len(rows))
	output := &Output{Message: rows}
	err = context.SetOutputObject(output)
	if err != nil {
		return false, err
	}

	logger.Info("binreader activity completed")
	return true, nil
}

func readNextBytes(file *os.File, number int) ([]byte, error) {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	return bytes, err
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
