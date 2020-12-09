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
		"(id, version, data_type, logical_type, name, physical_name, sequence, param_group_id, created_by, created_time, modified_by, modified_time) VALUES " +
		"(?,?,?,?,?,?,?,?,'admin',NOW(),'admin',NOW())"

	uid := uuid.New().String()
	_, err = tx.Exec(stmt1, uid, 0, "DOUBLE", "DEFAULT", paramName, paramName, 0, paramGroupID)

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

func CreateParamSpecWithModel(paramID, upper, target, lower string) (string, error) {
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

	stmt1 := "INSERT INTO parameter_value (id, asset_id, param_id, props) VALUES (?,?,?,?) "

	uid := uuid.New().String()
	_, err = tx.Exec(stmt1, uid, nil, paramID,
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

func CreateParamSpecWithAsset(assetID, tenantID, paramGroupID string) (string, error) {
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
	// uid := uuid.New().String()

	// stmt1 := "INSERT INTO parameter_value (id, asset_id, param_id, props) " +
	// 	"SELECT  uuid() as id, '" + assetID + "' as asset_id, p.id as param_id, v.props " +
	// 	"FROM parameter p, parameter_value v " +
	// 	"WHERE p.param_group_id='" + paramGroupID + "' " +
	// 	"AND p.id=v.param_id "

	stmt1 := "INSERT INTO parameter_value (id, asset_id, param_id, props) " +
		"SELECT  uuid() as id, '" + assetID + "' as asset_id, p.id as param_id, v.props " +
		"FROM param_group g, parameter p, parameter_value v " +
		"WHERE g.tenant_id='" + tenantID + "' " +
		"AND g.id=p.param_group_id " +
		"AND p.param_group_id='" + paramGroupID + "' " +
		"AND p.id=v.param_id " +
		"AND v.asset_id is null "

	_, err = tx.Exec(stmt1)

	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	log.Println("This parameter spec has been created.")
	return "", err
}
