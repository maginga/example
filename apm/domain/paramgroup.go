package domain

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type ParamGroup struct {
	Id   string
	Name string
}

func GetParamGroup(tenantID string) ([]ParamGroup, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id, name string
	rows, err := db.Query("SELECT id, name FROM param_group WHERE tenant_id=?", tenantID)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var paramGroupList = []ParamGroup{}

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println("ID: " + id + ", Name: " + name)
		paramGroupList = append(paramGroupList, ParamGroup{id, name})
	}

	return paramGroupList, err
}

func CreateParamGroup(tenantID string, typeID string, groupName string) (string, error) {
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

	stmt1 := "INSERT INTO param_group " +
		"(id, version, tenant_id, name, type_id, created_by, created_time) VALUES " +
		"(?,?,?,?,?,?,NOW()) "

	uid := uuid.New().String()
	_, err = tx.Exec(stmt1, uid, 0, tenantID, groupName, typeID, "admin")

	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	log.Println("This parameter group has been created.")

	return uid, err
}
