package domain

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func CreateSensor(assetID, deviceID, paramGroupID, duration, sensorName string) (string, error) {
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

	stmt1 := "INSERT INTO sensor " +
		"(id, version, asset_id, collecting, device_id, duration, name, physical_name, url, created_by, created_time, modified_by, modified_time) VALUES " +
		"(?,?,?,?,?,?,?,?,?,'admin',NOW(),'admin',NOW()) "

	uid := uuid.New().String()
	_, err = tx.Exec(stmt1, uid, 0, assetID, 1, deviceID, duration, sensorName, sensorName, "modbus://10.0.0.2:502")

	if err != nil {
		log.Panic(err)
	}

	stmt2 := "INSERT INTO sensor_param_group_join " +
		"(asset_id, param_group_id, sensor_id) VALUES " +
		"(?,?,?) "

	_, err = tx.Exec(stmt2, assetID, paramGroupID, uid)

	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	// log.Println("Sensor ID: " + uid)
	log.Println("This sensor has been created.")

	return sensorName, err
}
