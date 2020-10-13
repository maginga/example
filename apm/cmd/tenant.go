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
	"example/apm/naming"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tenantCmd represents the tenant command
var tenantCmd = &cobra.Command{
	Use:   "tenant",
	Short: "(02) Create a Tenant with Nest, Role and Authority.",
	Long: `Create a Tenant with Nest, Role and Authority. 
for example: 
	apm create tenant [Organization Code] [root catalog ID]
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)

		url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
		db, err := sql.Open("mysql", url)
		if err != nil {
			logger.Fatal(err)
		}
		defer db.Close()

		if len(args) <= 1 {
			var id string
			var name string
			rows, err := db.Query("SELECT id, name FROM catalog WHERE parent_id is null")
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

			logger.Println("")
			logger.Println("Select the catalog to connect with the tenant.")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			logger.Panic(err)
		}
		defer tx.Rollback()

		orgCode := args[0]
		rootCatalogID := args[1]

		seed := time.Now().UTC().UnixNano()
		nameGenerator := naming.NewNameGenerator(seed)
		name := nameGenerator.Generate()

		tenantID := uuid.New().String()
		tenantName := strings.ToUpper(orgCode + "_TENANT_" + name)
		stmt := "INSERT INTO tenant (id, catalog_id, tenant_id) VALUES (?,?,?) "
		_, err = tx.Exec(stmt, tenantID, rootCatalogID, tenantName)
		if err != nil {
			logger.Panic(err)
		}

		stmtNest := "INSERT INTO nest " +
			"(id, version, name, stream_spec, tenant_id, created_by, created_time) VALUES " +
			"(?,?,?,?,?,?,NOW()) "

		name = nameGenerator.Generate()
		nestName := strings.ToLower(strings.ReplaceAll(tenantName, "_", "-") + "-nest-" + name)
		nestID := nestName
		_, err = tx.Exec(stmtNest, nestID, 0, nestName,
			`{
				"id": "`+nestID+`",
				"modelJobs": [],
				"parameters": [],
				"storage": {
				  "score": "apm_score_`+strings.ReplaceAll(nestName, "-", "_")+`",
				  "trace": "apm_trace_`+strings.ReplaceAll(nestName, "-", "_")+`"
				},
				"topic": {
				  "alarmSpec": "apm-alarm-spec-`+nestName+`",
				  "assetSpec": "apm-asset-spec-`+nestName+`",
				  "modelSpec": "apm-model-spec-`+nestName+`",
				  "score": "apm-score-`+nestName+`",
				  "source": "apm-trace-`+nestName+`"
				}
			  }`,
			tenantID, "CLI")

		if err != nil {
			logger.Panic(err)
		}

		logger.Println("Tenant: " + tenantName + ", ID: " + tenantID)
		logger.Println("Nest Name(ID): " + nestName)
		logger.Println("The Tenant has been created.")

		// create authority & role
		stmt0 := "INSERT INTO authority " +
			"(id, item_id, item_type) " +
			"SELECT concat('authority_', descendant) as id, descendant as item_id, 'CATALOG' as item_type " +
			"FROM catalog_tree WHERE ancestor='" + rootCatalogID + "'"

		_, err = tx.Exec(stmt0)
		if err != nil {
			logger.Panic(err)
		}

		id := tenantName + "_ROLE_GENERAL_USER"
		stmt1 := "INSERT INTO roles " +
			"(id, version, tenant_id, role_desc, role_name, role_predefined, created_by, created_time) VALUES " +
			"(?,?,?,?,?,?,?,NOW()) "

		_, err = tx.Exec(stmt1, id, 0, tenantName, "General User", "General-User", 1, "CLI")
		if err != nil {
			logger.Panic(err)
		}

		stmt1Dir := "INSERT INTO role_directory " +
			"(directory_id, directory_name, directory_type, role_id, created_time) VALUES " +
			"(?,?,?,?,NOW()) "

		_, err = tx.Exec(stmt1Dir, "guest", "Guest", "USER", id)
		if err != nil {
			logger.Panic(err)
		}

		idx := 1
		stmtAuth1 := "INSERT INTO authority_target (id, authority_id, target_id, target_type, type) " +
			"SELECT concat(a.id, '" + fmt.Sprintf("%v", idx) + "') as id, a.id as authority_id, '" + id + "' as target_id, " +
			"'GROUP' as target_type, 'READ' as type " +
			"FROM catalog_tree c, authority a " +
			"WHERE a.item_id=c.descendant " +
			"AND c.ancestor='" + rootCatalogID + "'"

		_, err = tx.Exec(stmtAuth1)
		if err != nil {
			logger.Panic(err)
		}

		id = tenantName + "_ROLE_SYSTEM_ADMIN"
		stmt2 := "INSERT INTO roles " +
			"(id, version, tenant_id, role_desc, role_name, role_predefined, created_by, created_time) VALUES " +
			"(?,?,?,?,?,?,?,NOW()) "

		_, err = tx.Exec(stmt2, id, 0, tenantName, "System Administrator", "System-Admin", 1, "CLI")
		if err != nil {
			logger.Panic(err)
		}

		stmt2Dir := "INSERT INTO role_directory " +
			"(directory_id, directory_name, directory_type, role_id, created_time) VALUES " +
			"(?,?,?,?,NOW()) "

		_, err = tx.Exec(stmt2Dir, "admin", "Admin", "USER", id)
		if err != nil {
			logger.Panic(err)
		}

		stmt2Perm1 := "INSERT INTO role_perm_join " +
			"(role_id, perm_id) VALUES " +
			"(?,?) "

		_, err = tx.Exec(stmt2Perm1, id, "1001")
		if err != nil {
			logger.Panic(err)
		}
		stmt2Perm2 := "INSERT INTO role_perm_join " +
			"(role_id, perm_id) VALUES " +
			"(?,?) "

		_, err = tx.Exec(stmt2Perm2, id, "1002")
		if err != nil {
			logger.Panic(err)
		}

		idx++
		stmtAuth2 := "INSERT INTO authority_target (id, authority_id, target_id, target_type, type) " +
			"SELECT concat(a.id, '" + fmt.Sprintf("%v", idx) + "') as id, a.id as authority_id, '" + id + "' as target_id, " +
			"'GROUP' as target_type, 'ALL' as type " +
			"FROM catalog_tree c, authority a " +
			"WHERE a.item_id=c.descendant " +
			"AND c.ancestor='" + rootCatalogID + "'"

		_, err = tx.Exec(stmtAuth2)
		if err != nil {
			logger.Panic(err)
		}

		id = tenantName + "_ROLE_SYSTEM_ASSET_MANAGER"
		stmt3 := "INSERT INTO roles " +
			"(id, version, tenant_id, role_desc, role_name, role_predefined, created_by, created_time) VALUES " +
			"(?,?,?,?,?,?,?,NOW()) "

		_, err = tx.Exec(stmt3, id, 0, tenantName, "Data Manager", "Data-Manager", 1, "CLI")
		if err != nil {
			logger.Panic(err)
		}

		stmt3Dir := "INSERT INTO role_directory " +
			"(directory_id, directory_name, directory_type, role_id, created_time) VALUES " +
			"(?,?,?,?,NOW()) "

		_, err = tx.Exec(stmt3Dir, "metatron", "Metatron", "USER", id)
		if err != nil {
			logger.Panic(err)
		}

		stmt3Perm := "INSERT INTO role_perm_join " +
			"(role_id, perm_id) VALUES " +
			"(?,?) "

		_, err = tx.Exec(stmt3Perm, id, "1002")
		if err != nil {
			logger.Panic(err)
		}

		idx++
		stmtAuth3 := "INSERT INTO authority_target (id, authority_id, target_id, target_type, type) " +
			"SELECT concat(a.id, '" + fmt.Sprintf("%v", idx) + "') as id, a.id as authority_id, '" + id + "' as target_id, " +
			"'GROUP' as target_type, 'ALL' as type " +
			"FROM catalog_tree c, authority a " +
			"WHERE a.item_id=c.descendant " +
			"AND c.ancestor='" + rootCatalogID + "'"

		_, err = tx.Exec(stmtAuth3)
		if err != nil {
			logger.Panic(err)
		}

		id = tenantName + "_ROLE_TENANT_ADMIN"
		stmt4 := "INSERT INTO roles " +
			"(id, version, tenant_id, role_desc, role_name, role_predefined, created_by, created_time) VALUES " +
			"(?,?,?,?,?,?,?,NOW()) "

		_, err = tx.Exec(stmt4, id, 0, tenantName, "Tenant Administrator", "Tenant-Admin", 1, "CLI")
		if err != nil {
			logger.Panic(err)
		}

		stmt4Dir := "INSERT INTO role_directory " +
			"(directory_id, directory_name, directory_type, role_id, created_time) VALUES " +
			"(?,?,?,?,NOW()) "

		_, err = tx.Exec(stmt4Dir, "tenant", "Tenant", "USER", id)
		if err != nil {
			logger.Panic(err)
		}

		stmt4Perm1 := "INSERT INTO role_perm_join " +
			"(role_id, perm_id) VALUES " +
			"(?,?) "

		_, err = tx.Exec(stmt4Perm1, id, "1000")
		if err != nil {
			logger.Panic(err)
		}

		stmt4Perm2 := "INSERT INTO role_perm_join " +
			"(role_id, perm_id) VALUES " +
			"(?,?) "

		_, err = tx.Exec(stmt4Perm2, id, "1001")
		if err != nil {
			logger.Panic(err)
		}

		stmt4Perm3 := "INSERT INTO role_perm_join " +
			"(role_id, perm_id) VALUES " +
			"(?,?) "

		_, err = tx.Exec(stmt4Perm3, id, "1002")
		if err != nil {
			logger.Panic(err)
		}

		err = tx.Commit()
		if err != nil {
			logger.Panic(err)
		}

		logger.Println("The roles and authority of the tenant were granted.")
	},
}

func init() {
	createCmd.AddCommand(tenantCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tenantCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tenantCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
