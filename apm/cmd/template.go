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

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "(04) Create a Template.",
	Long: `Create a Template.
For example: 
	apm create template <--show tenant
	apm create template [Tenant ID] <--show type
	apm create template [Tenant ID] [Type Id] [Template Name] 
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
			var name string
			rows, err := db.Query("SELECT tenant_id as name FROM tenant ")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&name)
				if err != nil {
					log.Fatal(err)
				}
				logger.Println("Tenant ID: " + name)
			}

			if len(name) <= 0 {
				logger.Println("Tenant ID does not exist.")
			}

			logger.Println("---")
			return
		} else if len(args) <= 1 {
			var id, name, role string
			rows, err := db.Query("SELECT id, name, role FROM type WHERE tenant_id=? ", args[0])
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&id, &name, &role)
				if err != nil {
					log.Fatal(err)
				}
				logger.Println("Type ID: " + id + ", Name: " + name + ", Role: " + role)
			}

			if len(id) <= 0 {
				logger.Println("Type ID does not exist.")
			}

			logger.Println("---")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			logger.Panic(err)
		}
		defer tx.Rollback()

		tenantID := args[0]
		templateName := args[2]
		typeID := args[1]

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
			"ASSET", typeID, "CLI")
		if err != nil {
			logger.Panic(err)
		}

		err = tx.Commit()
		if err != nil {
			logger.Panic(err)
		}

		logger.Println("Template Id: " + templateID)
		logger.Println("The template has been created.")
	},
}

func init() {
	createCmd.AddCommand(templateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// templateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
