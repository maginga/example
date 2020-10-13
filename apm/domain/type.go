package domain

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type TypeTuple struct {
	Id   string
	Name string
}

func CreateType(tenantID string, typeName string, typeRole string) (string, error) {
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

	stmt0 := "INSERT INTO type " +
		"(id, version, tenant_id, name, role, sequence, created_by, created_time) VALUES " +
		"(?,?,?,?,?,?,?,NOW())"

	typeID := uuid.New().String()
	_, err = tx.Exec(stmt0, typeID, 0, tenantID, typeName, typeRole, 0, "admin")
	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Type Role: " + typeRole + ", Type Id: " + typeID)
	log.Println("The type has been created.")

	return typeID, err
}

func GetTypeList(tenantID, roleName string) ([]TypeTuple, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id, name string
	rows, err := db.Query("SELECT id, name FROM type WHERE tenant_id=? and role=?",
		tenantID, roleName)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var typeTupleList = []TypeTuple{}

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println("ID: " + id + ", Name: " + name)
		typeTupleList = append(typeTupleList, TypeTuple{id, name})
	}

	return typeTupleList, err
}
