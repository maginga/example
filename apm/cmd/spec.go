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

// specCmd represents the spec command
var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "(09) Create a parameter spec.",
	Long: `Create a parameter spec. 
For example: apm create spec [Asset ID] [Parameter ID] [Upper Value] [Target Value] [Lower Value]
`,
	Args: cobra.MinimumNArgs(5),
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

		assetID := args[0]
		paramID := args[1]
		upperValue := args[2]
		targetValue := args[3]
		lowerValue := args[4]

		stmt1 := "INSERT INTO parameter_value " +
			"(id, asset_id, param_id, props) VALUES " +
			"(?,?,?,?) "

		uid := uuid.New().String()
		_, err = tx.Exec(stmt1, uid, assetID, paramID,
			`{
			"type": "default",
			"lowerLimit": `+lowerValue+`,
			"targetValue": `+targetValue+`,
			"upperLimit": `+upperValue+`
		  }`)

		if err != nil {
			logger.Panic(err)
		}

		err = tx.Commit()
		if err != nil {
			logger.Panic(err)
		}

		logger.Println("This parameter spec has been created.")
	},
}

func init() {
	createCmd.AddCommand(specCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// specCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// specCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
