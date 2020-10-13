package domain

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func CreateDevice(tenantID, ipAddress, macAddress, modelNumber, serialNumber, deviceName string) (string, error) {
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

	var props = `{"mainPanel":"LOCAL CONTROL PANEL-21", "secondaryPanel":"WS-2 PANEL"}`

	stmt1 := "INSERT INTO device " +
		"(id, version, tenant_id, ip_addr, mac_addr, model_num, name, props, serial_num, created_by, created_time) VALUES " +
		"(?,?,?,?,?,?,?,?,?,?,NOW()) "

	uid := uuid.New().String()
	_, err = tx.Exec(stmt1, uid, 0, tenantID, ipAddress, macAddress, modelNumber, deviceName, props, serialNumber, "admin")

	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	log.Println("This device has been created.")

	return deviceName, err
}

type DeviceTuple struct {
	Id   string
	Name string
}

func GetDeviceList(tenantID string) ([]DeviceTuple, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id, name string
	rows, err := db.Query("SELECT id, name FROM device WHERE tenant_id=?", tenantID)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var list = []DeviceTuple{}

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println("ID: " + id + ", Name: " + name)
		list = append(list, DeviceTuple{id, name})
	}

	return list, err
}
