package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
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

	for i := 1; i <= 100; i++ {
		assetName := "Pump" + strconv.Itoa(i)
		sensorName := "S" + strconv.Itoa(i)
		update(db, assetName, sensorName, &config)
	}

	log.Println("complete.")
}

func update(db *sql.DB, assetName, sensorName string, config *Config) {
	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	props := `[
			{
			  "dataType": "String",
			  "defaultValue": "",
			  "description": "",
			  "inputType": "TEXT",
			  "isChange": false,
			  "isError": false,
			  "isNew": false,
			  "name": "Series Number",
			  "referenceType": "",
			  "seq": "1",
			  "value": "5478421570102"
			},
			{
			  "dataType": "String",
			  "defaultValue": "",
			  "description": "",
			  "inputType": "TEXT",
			  "isChange": false,
			  "isError": false,
			  "isNew": false,
			  "name": "Manufacturer",
			  "referenceType": "",
			  "seq": "2",
			  "value": "Manufacturer0102"
			},
			{
			  "dataType": "String",
			  "defaultValue": "",
			  "description": "",
			  "inputType": "TEXT",
			  "isChange": false,
			  "isError": false,
			  "isNew": false,
			  "name": "Frame",
			  "referenceType": "",
			  "seq": "3",
			  "value": "Frame0102"
			}
		  ]`

	assetID := uuid.New().String()
	stmt0 := "INSERT INTO asset " +
		"(id, version, catalog_id, image_url, name, physical_name, props, template_id, type_id, created_by, created_time) " +
		"SELECT '" + assetID + "' as id, 0 as version, '" + config.CatalogId + "' as catalog_id, 'apm://images/asset/poc_asset_01' as image_url, " +
		"'" + assetName + "' as name, " +
		"'" + assetName + "' as physical_name, " +
		"'" + props + "' as props, " +
		"id as template_id, type_id, 'admin' as created_by, NOW() as created_time " +
		"FROM asset_template " +
		"WHERE tenant_id='" + config.TenantId + "' AND id='" + config.TemplateId + "'"

	_, err = tx.Exec(stmt0)

	if err != nil {
		log.Panic(err)
	}

	stmt1 := "INSERT INTO asset_catalog_join (catalog_id, asset_id) VALUES (?,?) "
	_, err = tx.Exec(stmt1, config.CatalogId, assetID)

	if err != nil {
		log.Panic(err)
	}

	stmt2 := "INSERT INTO nest_egg (asset_id, nest_id) VALUES (?,?) "
	_, err = tx.Exec(stmt2, assetID, config.NestId)

	if err != nil {
		log.Panic(err)
	}

	// pvid := uuid.New().String()
	// stmt3 := "INSERT INTO parameter_value (id, asset_id, param_id, props) " +
	// 	"SELECT '" + pvid + "' as id, '" + assetID + "' as asset_id, v.param_id, v.props as props " +
	// 	"FROM param_group g, parameter p, parameter_value v " +
	// 	"WHERE g.tenant_id='" + config.TenantId + "'" +
	// 	"AND g.id=p.param_group_id " +
	// 	"AND p.id=v.param_id " +
	// 	"AND v.asset_id is null"
	// _, err = tx.Exec(stmt3)

	// if err != nil {
	// 	log.Panic(err)
	// }

	stmt3 := "INSERT INTO parameter_value (id, asset_id, param_id, props) " +
		"SELECT  uuid() as id, '" + assetID + "' as asset_id, p.id as param_id, v.props " +
		"FROM param_group g, parameter p, parameter_value v " +
		"WHERE g.tenant_id='" + config.TenantId + "' " +
		"AND g.id=p.param_group_id " +
		"AND p.param_group_id='" + config.ParamGroupId + "' " +
		"AND p.id=v.param_id " +
		"AND v.asset_id is null "

	_, err = tx.Exec(stmt3)

	if err != nil {
		log.Panic(err)
	}

	stmt4 := "INSERT INTO sensor " +
		"(id, version, asset_id, collecting, device_id, duration, name, physical_name, url, created_by, created_time, modified_by, modified_time) VALUES " +
		"(?,?,?,?,?,?,?,?,?,'admin',NOW(),'admin',NOW()) "

	sensorID := uuid.New().String()
	_, err = tx.Exec(stmt4, sensorID, 0, assetID, 1, config.DeviceId, "PT1S", sensorName, sensorName, "modbus://10.0.0.2:502")

	if err != nil {
		log.Panic(err)
	}

	stmt5 := "INSERT INTO sensor_param_group_join " +
		"(asset_id, param_group_id, sensor_id) VALUES " +
		"(?,?,?) "

	_, err = tx.Exec(stmt5, assetID, config.ParamGroupId, sensorID)

	if err != nil {
		log.Panic(err)
	}

	uid := uuid.New().String()
	stmt6 := "INSERT INTO sensor_status " +
		"(id, collecting_rate, sensor_id, status) VALUES " +
		"(?,?,?,?) "

	_, err = tx.Exec(stmt6, uid, 100, sensorID, 1)

	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}
}
