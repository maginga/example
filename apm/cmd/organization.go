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
	"example/apm/client"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger *log.Logger

// organizationCmd represents the organization command
var organizationCmd = &cobra.Command{
	Use:   "organization",
	Short: "(00) Create a Organization.",
	Long: `Create a Organization.
For example: 
	apm create organization [Organization Name]
`,
	Run: func(cmd *cobra.Command, args []string) {
		logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)

		url := fmt.Sprintf("%v", viper.Get("metadata.discovery-url"))
		db, err := sql.Open("mysql", url)
		if err != nil {
			logger.Fatal(err)
		}
		defer db.Close()

		if len(args) <= 0 {
			var id string
			var code string
			rows, err := db.Query("SELECT id, org_code FROM user_org ORDER BY created_time desc ")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&id, &code)
				if err != nil {
					log.Fatal(err)
				}
				logger.Println("ID: " + id + ", Code: " + code)
			}

			logger.Println("")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			logger.Panic(err)
		}
		defer tx.Rollback()

		orgName := strings.ToUpper(args[0]) + "_ORG"
		orgID := "ID_" + orgName

		stmt := "INSERT INTO user_org " +
			"(id, version, org_name, org_code, org_user_count, org_group_count, created_by, created_time) " +
			"VALUES " +
			"(?, ?, ?, ?, ?, ?, ?, NOW())"

		_, err = tx.Exec(stmt, orgID, 1, orgName, orgName, 0, 0, "CLI")
		if err != nil {
			logger.Panic(err)
		}

		ugcStmt := "INSERT INTO user_org_member " +
			"(member_id, member_name, member_type, org_id) " +
			"SELECT id as member_id, group_name as member_name, 'GROUP' as member_type, '" + orgID + "' as org_id " +
			"FROM user_group "
		userGroupCount, err := tx.Exec(ugcStmt)
		if err != nil {
			logger.Panic(err)
		}

		ucStmt := "INSERT INTO user_org_member " +
			"(member_id, member_name, member_type, org_id) " +
			"SELECT user_name as member_id, user_full_name as member_name, 'USER' as member_type, '" + orgID + "' as org_id " +
			"FROM users WHERE user_name in ('admin', 'guest', 'polaris') "

		userCount, err := tx.Exec(ucStmt)
		if err != nil {
			logger.Panic(err)
		}

		uc, _ := userCount.RowsAffected()
		gc, _ := userGroupCount.RowsAffected()
		updateSQL := "UPDATE user_org SET org_user_count=" + fmt.Sprintf("%v", uc) + ", org_group_count=" + fmt.Sprintf("%v", gc) + " WHERE id='" + orgID + "'"
		_, err = tx.Exec(updateSQL)

		if err != nil {
			logger.Panic(err)
		}

		err = tx.Commit()
		if err != nil {
			logger.Panic(err)
		}

		logger.Println("Organization Code(Tenant): " + orgName)
		logger.Println("The organization of the company has been created.")

		discoveryURL := fmt.Sprintf("%v", viper.Get("discovery.url"))
		logger.Println("Discovery URL: " + discoveryURL)
		c, err := client.NewDruidClient(discoveryURL)
		if err != nil {
			panic(err)
		}

		t, err := c.InitOrganization()
		if err != nil {
			panic(err)
		}

		logger.Println("The organization was initialized.: " + t)

		q, err := c.ConfigureUser()
		if err != nil {
			panic(err)
		}

		logger.Println("Login Setup Completed.: " + q)
	},
}

func init() {
	createCmd.AddCommand(organizationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// organizationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// organizationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
