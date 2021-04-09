package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Parser struct {
	header      []string
	current     bool
	temperature bool
	vibration   bool
	rows        [][]string
}

func (p *Parser) WriteFile(fileName string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(strings.Join(p.header, ",") + "\n"); err != nil {
		log.Println(err)
	}

	window := make(map[string]interface{})

	for _, row := range p.rows {
		time := row[0]
		assetName := row[3]
		sensor := row[4]
		paramType := row[5]
		length, _ := strconv.Atoi(row[6])

		val := make([]string, 0, length)
		for i := 0; i < length; i++ {
			val = append(val, row[i+7])
		}
		window[paramType] = val

		if paramType == "CUR" {
			p.current = true
		} else if paramType == "TEM" {
			p.temperature = true
		} else {
			p.vibration = true
		}

		if p.current && p.temperature && p.vibration {
			p.current = false
			p.temperature = false
			p.vibration = false

			result := time + "," + assetName + "," + sensor + "," +
				strings.Join(window["CUR"].([]string), ",") + "," +
				strings.Join(window["TEM"].([]string), ",") + "," +
				strings.Join(window["VIR"].([]string), ",") + "\n"

			if _, err := f.WriteString(result); err != nil {
				log.Println(err)
			}
		}
	}
}

func (p *Parser) ReadFile(fileName string) {
	f, _ := filepath.Abs(fileName)
	csvfile, err := os.Open(f)
	if err != nil {
		log.Println(err)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	//reader.FieldsPerRecord = -1
	p.rows, _ = reader.ReadAll()
}

func NewParser() (*Parser, error) {
	header := []string{"event_time", "assetName", "sensorLocation",
		"Current1", "Current2", "Current3", "Current4", "Current5", "Current6",
		"Temperature1", "Temperature2",
		"VibrationX", "VibrationY", "VibrationZ", "VibrationXX", "VibrationYY"}

	parser := &Parser{header, false, false, false, nil}
	return parser, nil
}
