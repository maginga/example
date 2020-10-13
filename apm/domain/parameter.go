package domain

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func CreateParameter(paramGroupID string, paramName string) (string, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	stmt1 := "INSERT INTO parameter " +
		"(id, version, data_type, logical_type, name, physical_name, sequence, param_group_id, created_by, created_time) VALUES " +
		"(?,?,?,?,?,?,?,?,?,NOW()) "

	uid := uuid.New().String()
	_, err = tx.Exec(stmt1, uid, 0, "DOUBLE", "DEFAULT", paramName, paramName, 0, paramGroupID, "admin")

	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	log.Println("This parameter has been created.")
	return uid, err
}

func CreateParameterWithSpec(assetID, paramID, upper, target, lower string) (string, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	stmt1 := "INSERT INTO parameter_value " +
		"(id, asset_id, param_id, props) VALUES " +
		"(?,?,?,?) "

	uid := uuid.New().String()
	_, err = tx.Exec(stmt1, uid, assetID, paramID,
		`{
			"type": "default",
			"lowerLimit": `+lower+`,
			"targetValue": `+target+`,
			"upperLimit": `+upper+`
		  }`)

	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	log.Println("This parameter spec has been created.")
	return uid, err
}
