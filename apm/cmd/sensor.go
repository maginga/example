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

// sensorCmd represents the sensor command
var sensorCmd = &cobra.Command{
	Use:   "sensor",
	Short: "(08) Create a Sensor of Asset and link to the ParameterGroup.",
	Long: `Create a Sensor of Asset and link to the ParameterGroup. 
For example: 
	apm create sensor [Asset ID] [Sensor Name] [Duration(PT5S)] [Param Group ID]
`,
	Args: cobra.MinimumNArgs(3),
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
		sensorName := args[1]
		duration := args[2]
		paramGroupID := args[3]

		stmt1 := "INSERT INTO sensor " +
			"(id, version, asset_id, collecting, duration, name, physical_name, url, created_by, created_time) VALUES " +
			"(?,?,?,?,?,?,?,?,?,NOW()) "

		uid := uuid.New().String()
		_, err = tx.Exec(stmt1, uid, 0, assetID, 1, duration, sensorName, sensorName, "modbus://10.0.0.2:502", "CLI")

		if err != nil {
			logger.Panic(err)
		}

		stmt2 := "INSERT INTO sensor_param_group_join " +
			"(asset_id, param_group_id, sensor_id) VALUES " +
			"(?,?,?) "

		_, err = tx.Exec(stmt2, assetID, paramGroupID, uid)

		if err != nil {
			logger.Panic(err)
		}

		err = tx.Commit()
		if err != nil {
			logger.Panic(err)
		}

		logger.Println("Sensor ID: " + uid)
		logger.Println("This sensor has been created.")
	},
}

func init() {
	createCmd.AddCommand(sensorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sensorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sensorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
