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

// paramgroupCmd represents the paramgroup command
var paramgroupCmd = &cobra.Command{
	Use:   "paramgroup",
	Short: "(05) Create a parameter group.",
	Long: `Create a parameter group.
For example: 
	apm create paramgroup <--show tenant / type
	apm create paramgroup [Group Name] [Tenant ID] [Type ID]
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
			var id, tenant, name string
			rows, err := db.Query("SELECT id, tenant_id, name FROM type WHERE role='PARAM' ")
			if err != nil {
				logger.Panic(err)
			}
			defer rows.Close()

			logger.Println("---")
			for rows.Next() {
				err := rows.Scan(&id, &tenant, &name)
				if err != nil {
					log.Fatal(err)
				}
				logger.Println("Tenant: " + tenant + ", Type ID: " + id + ", Type Name: " + name)
				logger.Println("---")
			}

			if len(name) <= 0 {
				logger.Println("Parameter Group does not exist.")
			}
			return
		}

		tx, err := db.Begin()
		if err != nil {
			logger.Panic(err)
		}
		defer tx.Rollback()

		groupName := args[0]
		tenantID := args[1]
		typeID := args[2]

		stmt1 := "INSERT INTO param_group " +
			"(id, version, tenant_id, name, type_id, created_by, created_time) VALUES " +
			"(?,?,?,?,?,?,NOW()) "

		uid := uuid.New().String()
		_, err = tx.Exec(stmt1, uid, 0, tenantID, groupName, typeID, "CLI")

		if err != nil {
			logger.Panic(err)
		}

		err = tx.Commit()
		if err != nil {
			logger.Panic(err)
		}

		logger.Println("Parameter Group ID: " + uid)
		logger.Println("This parameter group has been created.")
	},
}

func init() {
	createCmd.AddCommand(paramgroupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// paramgroupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// paramgroupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
