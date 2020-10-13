package domain

import (
	"database/sql"
	"example/apm/naming"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type NestTuple struct {
	Id   string
	Name string
}

func GetNestList(tenantID string) ([]NestTuple, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id, name string
	rows, err := db.Query("SELECT id, name FROM nest WHERE tenant_id in (select id from tenant where tenant_id=?)", tenantID)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var list = []NestTuple{}

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println("ID: " + id + ", Name: " + name)
		list = append(list, NestTuple{id, name})
	}

	return list, err
}

type TenantTuple struct {
	Id   string
	Name string
}

func GetTenantList() ([]TenantTuple, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id, name string
	rows, err := db.Query("SELECT tenant_id, tenant_id FROM tenant")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var list = []TenantTuple{}

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println("ID: " + id + ", Name: " + name)
		list = append(list, TenantTuple{id, name})
	}

	return list, err
}

func CreateTenant(orgCode string, rootCatalogID string) (string, string, string, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	seed := time.Now().UTC().UnixNano()
	nameGenerator := naming.NewNameGenerator(seed)
	name := nameGenerator.Generate()

	tenantID := uuid.New().String()
	tenantName := orgCode
	stmt := "INSERT INTO tenant (id, catalog_id, tenant_id) VALUES (?,?,?) "
	_, err = tx.Exec(stmt, tenantID, rootCatalogID, tenantName)
	if err != nil {
		log.Panic(err)
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
		tenantID, "admin")

	if err != nil {
		log.Panic(err)
	}

	log.Println("The Tenant has been created.")

	// stmt0 := "INSERT INTO authority " +
	// 	"(id, item_id, item_type) " +
	// 	"SELECT concat('authority_', descendant) as id, descendant as item_id, 'CATALOG' as item_type " +
	// 	"FROM catalog_tree WHERE ancestor='" + rootCatalogID + "'"

	// _, err = tx.Exec(stmt0)
	// if err != nil {
	// 	log.Panic(err)
	// }

	id := tenantName + "_ROLE_GENERAL_USER"
	stmt1 := "INSERT INTO roles " +
		"(id, version, tenant_id, role_desc, role_name, role_predefined, created_by, created_time) VALUES " +
		"(?,?,?,?,?,?,?,NOW()) "

	_, err = tx.Exec(stmt1, id, 0, tenantName, "General User", "["+id+"]General-User", 1, "admin")
	if err != nil {
		log.Panic(err)
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
		log.Panic(err)
	}

	id = tenantName + "_ROLE_SYSTEM_ADMIN"
	stmt2 := "INSERT INTO roles " +
		"(id, version, tenant_id, role_desc, role_name, role_predefined, created_by, created_time) VALUES " +
		"(?,?,?,?,?,?,?,NOW()) "

	_, err = tx.Exec(stmt2, id, 0, tenantName, "System Administrator", "["+id+"]System-Admin", 1, "admin")
	if err != nil {
		log.Panic(err)
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
		log.Panic(err)
	}

	id = tenantName + "_ROLE_SYSTEM_ASSET_MANAGER"
	stmt3 := "INSERT INTO roles " +
		"(id, version, tenant_id, role_desc, role_name, role_predefined, created_by, created_time) VALUES " +
		"(?,?,?,?,?,?,?,NOW()) "

	_, err = tx.Exec(stmt3, id, 0, tenantName, "Data Manager", "["+id+"]Data-Manager", 1, "admin")
	if err != nil {
		log.Panic(err)
	}

	stmt3Perm := "INSERT INTO role_perm_join (role_id, perm_id) VALUES (?,?) "
	_, err = tx.Exec(stmt3Perm, id, "1002")
	if err != nil {
		log.Panic(err)
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
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	log.Println("The role was created for each tenant.")
	return tenantName, tenantName, nestName, err
}

func CreateRoleOfTenant(orgCode, rootCatalogID, userID string) error {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	tenantName := orgCode
	id := tenantName + "_ROLE_SYSTEM_ASSET_MANAGER"

	stmt3Dir := "INSERT INTO role_directory " +
		"(directory_id, directory_name, directory_type, role_id, created_time) VALUES " +
		"(?,?,?,?,NOW()) "

	_, err = tx.Exec(stmt3Dir, userID, userID, "USER", id)
	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	log.Println("The role directory was created.")
	return err
}
