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
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// typeCmd represents the type command
var typeCmd = &cobra.Command{
	Use:   "type",
	Short: "(03) Create a Type. (ASSET, PARAM)",
	Long: `Create a Type (ASSET, PARAM).
For example: 
	apm create type [Tenant ID] [Type Role(ASSET, PARAM)] [Type Name]
`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)

		url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
		db, err := sql.Open("mysql", url)
		if err != nil {
			logger.Fatal(err)
		}
		defer db.Close()

		if len(args) <= 2 {
			var id, name string
			rows, err := db.Query("SELECT id, name FROM type WHERE tenant_id=? and role=? ", args[0], args[1])
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&id, &name)
				if err != nil {
					log.Fatal(err)
				}
				logger.Println("ID: " + id + ", Name: " + name)
			}

			if len(id) <= 0 {
				logger.Println("Type does not exist.")
			}

			logger.Println("---")
			return
		}

		tenantID := args[0]
		typeRole := strings.ToUpper(args[1])
		typeName := args[2]

		var typeID string
		sql := "SELECT id FROM type WHERE tenant_id='" + tenantID +
			"' AND role='" + typeRole + "' AND name='" + typeName + "'"
		err = db.QueryRow(sql).Scan(&typeID)
		if err != nil {
			logger.Println("Type does not exist: " + err.Error())

			tx, err := db.Begin()
			if err != nil {
				logger.Panic(err)
			}
			defer tx.Rollback()

			stmt0 := "INSERT INTO type " +
				"(id, version, tenant_id, name, role, sequence, created_by, created_time) VALUES " +
				"(?,?,?,?,?,?,?,NOW())"

			typeID = uuid.New().String()
			_, err = tx.Exec(stmt0, typeID, 0, tenantID, typeName, typeRole, 0, "CLI")
			if err != nil {
				logger.Panic(err)
			}

			err = tx.Commit()
			if err != nil {
				logger.Panic(err)
			}

			logger.Println("Type Role: " + typeRole + ", Type Id: " + typeID)
			logger.Println("The type has been created.")
		} else {
			logger.Println("Existed Type Role: " + typeRole + ", Type Id: " + typeID)
		}
	},
}

func init() {
	createCmd.AddCommand(typeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// typeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// typeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
