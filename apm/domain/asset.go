package domain

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func CreateAsset(tenantID, templateID, catalogID, nestID, assetName string) (string, string, error) {
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

	uid := uuid.New().String()
	stmt0 := "INSERT INTO asset " +
		"(id, version, catalog_id, image_url, name, physical_name, props, template_id, type_id, created_by, created_time) " +
		"SELECT '" + uid + "' as id, 0 as version, '" + catalogID + "' as catalog_id, 'apm://images/asset/poc_asset_01' as image_url, " +
		"'" + assetName + "' as name, " +
		"'" + assetName + "' as physical_name, " +
		"'" + props + "' as props, " +
		"id as template_id, type_id, 'admin' as created_by, NOW() as created_time " +
		"FROM asset_template " +
		"WHERE tenant_id='" + tenantID + "' AND id='" + templateID + "'"

	_, err = tx.Exec(stmt0)

	if err != nil {
		log.Panic(err)
	}

	stmt1 := "INSERT INTO asset_catalog_join (catalog_id, asset_id) VALUES (?,?) "
	_, err = tx.Exec(stmt1, catalogID, uid)

	if err != nil {
		log.Panic(err)
	}

	stmt2 := "INSERT INTO nest_egg (asset_id, nest_id) VALUES (?,?) "
	_, err = tx.Exec(stmt2, uid, nestID)

	if err != nil {
		log.Panic(err)
	}

	// pvid := uuid.New().String()
	// stmt3 := "INSERT INTO parameter_value (id, asset_id, param_id, props) " +
	// 	"SELECT '" + pvid + "' as id, '" + uid + "' as asset_id, v.param_id, v.props as props " +
	// 	"FROM param_group g, parameter p, parameter_value v " +
	// 	"WHERE g.tenant_id='" + tenantID + "'" +
	// 	"AND g.id=p.param_group_id " +
	// 	"AND p.id=v.param_id " +
	// 	"AND v.asset_id is null"
	// _, err = tx.Exec(stmt3)

	// if err != nil {
	// 	log.Panic(err)
	// }

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	// log.Println("Asset ID: " + uid + ", Physical Name: " + assetName)
	log.Println("This asset has connected to the Nest.")
	return assetName, uid, err
}
