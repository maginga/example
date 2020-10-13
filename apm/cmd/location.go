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

// locationCmd represents the location command
var locationCmd = &cobra.Command{
	Use:   "location",
	Short: "(01) Create a Location.",
	Long: `Create a Location. 
for example: 
	apm create location [root]
	apm create location [child] [parents] [depth]
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)

		url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
		logger.Println("URL: " + url)
		db, err := sql.Open("mysql", url)
		if err != nil {
			logger.Fatal(err)
		}
		defer db.Close()

		tx, err := db.Begin()
		if err != nil {
			logger.Panic(err)
		}
		defer tx.Rollback()

		leafName := args[0]
		id := uuid.New().String()

		if len(args) <= 1 {
			stmt1 := "INSERT INTO catalog " +
				"(id, version, name, locking, sequence, created_by, created_time) VALUES " +
				"(?,?,?,?,?,?,NOW()) "

			_, err = tx.Exec(stmt1, id, 0, leafName, 0, 1, "CLI")
			if err != nil {
				logger.Panic(err)
			}

			stmt2 := "INSERT INTO catalog_tree " +
				"(ancestor, descendant, depth) VALUES " +
				"(?,?,?) "

			_, err = tx.Exec(stmt2, id, id, 0)
			if err != nil {
				logger.Panic(err)
			}
			logger.Println("Root Catalog ID: " + id)
		} else {
			parents := args[1]
			depth := args[2]

			stmt1 := "INSERT INTO catalog (id, version, name, locking, sequence, parent_id, created_by, created_time) " +
				"SELECT '" + id + "' as id, 0 as version, '" + leafName + "' as name, 0 as locking, 1 as sequence, " +
				"id as parent_id, 'CLI' as created_by, NOW() as created_time " +
				"FROM catalog WHERE name='" + parents + "'"

			_, err = tx.Exec(stmt1)
			if err != nil {
				logger.Panic(err)
			}

			stmt2 := "INSERT INTO catalog_tree " +
				"SELECT id as ancestor, '" + id + "' as descendant, " + fmt.Sprintf("%v", depth) +
				" as depth FROM catalog WHERE name='" + parents + "'"

			_, err = tx.Exec(stmt2)
			if err != nil {
				logger.Panic(err)
			}
		}

		err = tx.Commit()
		if err != nil {
			logger.Panic(err)
		}

		logger.Println("The location was created.")
	},
}

func init() {
	createCmd.AddCommand(locationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// locationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// locationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
