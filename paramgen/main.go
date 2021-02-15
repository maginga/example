package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gonum/stat"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

func main() {
	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql", config.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fileName, _ := filepath.Abs(config.FilePath)
	csvfile, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	//reader.FieldsPerRecord = -1

	rows, _ := reader.ReadAll()

	colCnt := len(rows[0])
	rowCnt := len(rows)

	for c := 2; c < colCnt; c++ {
		colValues := make([]float64, 0)
		columnName := rows[0][c]

		for r := 1; r < rowCnt; r++ {
			strVal := rows[r][c]
			f, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			colValues = append(colValues, f)
		}

		if len(colValues) > 1 {
			lower, upper, target := findSpec(colValues)
			update(db, config.ParameterGroup, columnName, lower, upper, target)
		}
	}
}

func update(db *sql.DB, group, param string, l, u, t float64) {
	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	stmt1 := "INSERT INTO parameter " +
		"(id, version, data_type, logical_type, name, physical_name, sequence, param_group_id, created_by, created_time, modified_by, modified_time) VALUES " +
		"(?,?,?,?,?,?,?,?,'admin',NOW(),'admin',NOW())"

	paramID := uuid.New().String()
	_, err = tx.Exec(stmt1, paramID, 0, "DOUBLE", "DEFAULT", param, param, 0, group)

	if err != nil {
		log.Panic(err)
	}

	stmt2 := "INSERT INTO parameter_value (id, asset_id, param_id, props) VALUES (?,?,?,?) "

	pvID := uuid.New().String()
	_, err = tx.Exec(stmt2, pvID, nil, paramID,
		`{
			"type": "default",
			"lowerLimit": `+fmt.Sprintf("%f", l)+`,
			"targetValue": `+fmt.Sprintf("%f", t)+`,
			"upperLimit": `+fmt.Sprintf("%f", u)+`
		  }`)

	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}
}

func findSpec(a []float64) (lo, up, avg float64) {
	mean := stat.Mean(a, nil)
	stdev := stat.StdDev(a, nil)

	lo = mean - (3 * stdev)
	up = mean + (3 * stdev)

	return lo, up, mean
}