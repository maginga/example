/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// assetCmd represents the asset command
var assetCmd = &cobra.Command{
	Use:   "asset",
	Short: "(07) Create a asset.",
	Long: `Create a asset.
For example: 
	apm create asset <--show tenant, nest, template
	apm create asset [Root Catalog ID] <--show catalogs
	apm create asset [Template ID] [Catalog ID] [Tenant ID] [Nest ID] [Asset Name]
`,
	Run: func(cmd *cobra.Command, args []string) {
		logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)

		url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
		db, err := sql.Open("mysql", url)
		if err != nil {
			logger.Fatal(err)
		}
		defer db.Close()

		if len(args) <= 0 {
			var id, tenant, name, nest, catalog, pgid, pgname string
			rows, err := db.Query("SELECT a.id, a.tenant_id as tenant, a.name, n.id as nest, t.catalog_id catalog, " +
				"p.id pgid, p.name pgname " +
				"FROM asset_template a, tenant t, nest n, param_group p " +
				"WHERE t.tenant_id=a.tenant_id AND t.id=n.tenant_id AND a.role='ASSET' and t.tenant_id=p.tenant_id ")
			if err != nil {
				logger.Panic(err)
			}
			defer rows.Close()

			logger.Println("---")
			for rows.Next() {
				err := rows.Scan(&id, &tenant, &name, &nest, &catalog, &pgid, &pgname)
				if err != nil {
					log.Fatal(err)
				}
				logger.Println("Template ID: " + id + ", Template Name: " + name)
				logger.Println("Tenant: " + tenant + ", Catalog: " + catalog)
				logger.Println("Nest: " + nest)
				logger.Println("Param Group ID: " + pgid + ", Param Group Name: " + pgname)
				logger.Println("---")
			}

			if len(name) <= 0 {
				logger.Println("meta does not exist.")
			}
			return
		}

		if len(args) <= 1 {
			var id, name string
			rows, err := db.Query("WITH RECURSIVE tree as ( " +
				"SELECT id, name FROM catalog WHERE id='" + args[0] + "' " +
				"UNION " +
				"SELECT catalog.id, catalog.name FROM catalog, tree WHERE tree.id=catalog.parent_id " +
				") " +
				"SELECT * FROM tree")

			if err != nil {
				logger.Panic(err)
			}
			defer rows.Close()

			logger.Println("---")
			for rows.Next() {
				err := rows.Scan(&id, &name)
				if err != nil {
					log.Fatal(err)
				}
				logger.Println("Catalog ID: " + id + ", Catalog Name: " + name)
				logger.Println("---")
			}

			if len(name) <= 0 {
				logger.Println("Catalog does not exist.")
			}
			return
		}

		tx, err := db.Begin()
		if err != nil {
			logger.Panic(err)
		}
		defer tx.Rollback()

		templateID := args[0]
		catalogID := args[1]
		tenantID := args[2]
		nestID := args[3]
		assetName := args[4]

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
			"(id, version, catalog_id, name, physical_name, props, template_id, type_id, created_by, created_time) " +
			"SELECT '" + uid + "' as id, 0 as version, '" + catalogID + "' as catalog_id, '" + assetName + "' as name, " +
			"'" + assetName + "' as physical_name, " +
			"'" + props + "' as props, " +
			"id as template_id, type_id, 'CLI' as created_by, NOW() as created_time " +
			"FROM asset_template " +
			"WHERE tenant_id='" + tenantID + "' AND id='" + templateID + "'"

		_, err = tx.Exec(stmt0)

		// stmt0 := "INSERT INTO asset " +
		// 	"(id, version, catalog_id, name, physical_name, props, template_id, type_id, created_by, created_time) VALUES " +
		// 	"(?,?,?,?,?,?,?,?,?,NOW()) "
		// _, err = tx.Exec(stmt0, uid, 0, catalogID, assetName, assetName,
		// 	`[
		// 		{
		// 		  "dataType": "String",
		// 		  "defaultValue": "",
		// 		  "description": "",
		// 		  "inputType": "TEXT",
		// 		  "isChange": false,
		// 		  "isError": false,
		// 		  "isNew": false,
		// 		  "name": "Series Number",
		// 		  "referenceType": "",
		// 		  "seq": "1",
		// 		  "value": "5478421570102"
		// 		},
		// 		{
		// 		  "dataType": "String",
		// 		  "defaultValue": "",
		// 		  "description": "",
		// 		  "inputType": "TEXT",
		// 		  "isChange": false,
		// 		  "isError": false,
		// 		  "isNew": false,
		// 		  "name": "Manufacturer",
		// 		  "referenceType": "",
		// 		  "seq": "2",
		// 		  "value": "Manufacturer0102"
		// 		},
		// 		{
		// 		  "dataType": "String",
		// 		  "defaultValue": "",
		// 		  "description": "",
		// 		  "inputType": "TEXT",
		// 		  "isChange": false,
		// 		  "isError": false,
		// 		  "isNew": false,
		// 		  "name": "Frame",
		// 		  "referenceType": "",
		// 		  "seq": "3",
		// 		  "value": "Frame0102"
		// 		}
		// 	  ]`,
		// 	templateID,
		// 	typeID, "CLI")

		if err != nil {
			logger.Panic(err)
		}

		stmt1 := "INSERT INTO asset_catalog_join (catalog_id, asset_id) VALUES (?,?) "
		_, err = tx.Exec(stmt1, catalogID, uid)

		if err != nil {
			logger.Panic(err)
		}

		stmt2 := "INSERT INTO nest_egg (asset_id, nest_id) VALUES (?,?) "
		_, err = tx.Exec(stmt2, uid, nestID)

		if err != nil {
			logger.Panic(err)
		}

		err = tx.Commit()
		if err != nil {
			logger.Panic(err)
		}

		logger.Println("Asset ID: " + uid + ", Physical Name: " + assetName)
		logger.Println("This asset has connected to the Nest.")
	},
}

func init() {
	createCmd.AddCommand(assetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// assetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// assetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
