package domain

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type OrgTuple struct {
	Id   string
	Name string
}

func GetOrgList() ([]OrgTuple, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.discovery-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id, name string
	rows, err := db.Query("select id, org_code from user_org")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var list = []OrgTuple{}

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println("ID: " + id + ", Name: " + name)
		list = append(list, OrgTuple{id, name})
	}

	return list, err
}

func CreateOrganization(name string) (string, string, error) {
	// logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
	url := fmt.Sprintf("%v", viper.Get("metadata.discovery-url"))
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

	orgName := strings.ToUpper(name) + "_ORG"
	orgID := "ID_" + orgName

	stmt := "INSERT INTO user_org " +
		"(id, version, org_name, org_code, org_user_count, org_group_count, created_by, created_time) " +
		"VALUES " +
		"(?, ?, ?, ?, ?, ?, ?, NOW())"

	_, err = tx.Exec(stmt, orgID, 1, orgName, orgName, 0, 0, "admin")
	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	// log.Println("Organization Code(Tenant): " + orgName)
	log.Println("The organization has been created.")

	return orgName, orgID, err
}

type UserTuple struct {
	Id   string
	Name string
}

func GetUserList() ([]UserTuple, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.discovery-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id, name string
	rows, err := db.Query("select user_name, user_full_name from users where id not in (select member_id from user_org_member)")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var list = []UserTuple{}

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println("ID: " + id + ", Name: " + name)
		list = append(list, UserTuple{id, name})
	}

	return list, err
}

func GetUser(tenantID string) ([]UserTuple, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.discovery-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id, name string
	rows, err := db.Query("select member_id, member_name from user_org_member "+
		"where member_type='USER' and org_id in (select id from user_org where org_code=?)", tenantID)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var list = []UserTuple{}

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, UserTuple{id, name})
	}

	return list, err
}

func LinkOrgMember(orgID, userID, userName string) error {
	url := fmt.Sprintf("%v", viper.Get("metadata.discovery-url"))
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

	ucStmt := "INSERT INTO user_org_member (member_id, member_name, member_type, org_id) VALUES (?,?,?,?)"
	userCount, err := tx.Exec(ucStmt, userID, userName, "USER", orgID)
	if err != nil {
		log.Panic(err)
	}

	uc, _ := userCount.RowsAffected()
	gc := 0
	updateSQL := "UPDATE user_org SET org_user_count=" + fmt.Sprintf("%v", uc) + ", org_group_count=" + fmt.Sprintf("%v", gc) + " WHERE id='" + orgID + "'"
	_, err = tx.Exec(updateSQL)

	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	// log.Println("Organization Code(Tenant): " + orgName)
	log.Println("You have connected users to your organization.")

	return err
}
