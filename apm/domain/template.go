package domain

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type TemplateTuple struct {
	Id   string
	Name string
}

func GetTemplateList(tenantID string) ([]TemplateTuple, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id, name string
	rows, err := db.Query("SELECT id, name FROM asset_template WHERE tenant_id=?", tenantID)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var list = []TemplateTuple{}

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println("ID: " + id + ", Name: " + name)
		list = append(list, TemplateTuple{id, name})
	}

	return list, err
}

func CreateTemplate(tenantID string, templateName string, typeID string) (string, error) {
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

	stmtTemplate := "INSERT INTO asset_template " +
		"(id, version, tenant_id, name, image_url, props, role, type_id, created_by, created_time) VALUES " +
		"(?,?,?,?,?,?,?,?,?,NOW()) "

	templateID := uuid.New().String()
	_, err = tx.Exec(stmtTemplate,
		templateID, 0, tenantID, templateName, "apm://images/asset/poc_asset_model_01",
		`[
				{
				  "dataType": "String",
				  "defaultValue": "",
				  "description": "",
				  "inputType": "TEXT",
				  "name": "Series Number",
				  "referenceType": "",
				  "seq": "1"
				},
				{
				  "dataType": "String",
				  "defaultValue": "",
				  "description": "",
				  "inputType": "TEXT",
				  "name": "Manufacturer",
				  "referenceType": "",
				  "seq": "2"
				},
				{
				  "dataType": "String",
				  "defaultValue": "",
				  "description": "",
				  "inputType": "TEXT",
				  "name": "Frame",
				  "referenceType": "",
				  "seq": "3"
				}
			  ]`,
		"ASSET", typeID, "admin")
	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	log.Println("The template has been created.")

	return templateID, err
}
