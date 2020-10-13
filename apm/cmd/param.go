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

// paramCmd represents the param command
var paramCmd = &cobra.Command{
	Use:   "param",
	Short: "(06) Create a parameter into parameter group.",
	Long: `Create a parameter into parameter group.
For example: 
	apm create param [Parameter Name] [Param Group ID]
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

		tx, err := db.Begin()
		if err != nil {
			logger.Panic(err)
		}
		defer tx.Rollback()

		paramName := args[0]
		paramGroupID := args[1]

		stmt1 := "INSERT INTO parameter " +
			"(id, version, data_type, logical_type, name, physical_name, sequence, param_group_id, created_by, created_time) VALUES " +
			"(?,?,?,?,?,?,?,?,?,NOW()) "

		uid := uuid.New().String()
		_, err = tx.Exec(stmt1, uid, 0, "DOUBLE", "DEFAULT", paramName, paramName, 0, paramGroupID, "CLI")

		if err != nil {
			logger.Panic(err)
		}

		err = tx.Commit()
		if err != nil {
			logger.Panic(err)
		}

		logger.Println("Group ID: " + paramGroupID + ", Parameter ID: " + uid)
		logger.Println("This parameter has been created.")
	},
}

func init() {
	createCmd.AddCommand(paramCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// paramCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// paramCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
