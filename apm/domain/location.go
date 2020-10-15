package domain

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func CreateLocation(child string, parents string, depth string) error {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	log.Println("URL: " + url)
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

	id := uuid.New().String()

	if child == parents {
		stmt1 := "INSERT INTO catalog " +
			"(id, version, name, locking, sequence, created_by, created_time) VALUES " +
			"(?,?,?,?,?,?,NOW()) "

		_, err = tx.Exec(stmt1, id, 0, child, 0, 1, "admin")
		if err != nil {
			log.Panic(err)
		}

		// stmt2 := "INSERT INTO catalog_tree " +
		// 	"(ancestor, descendant, depth) VALUES " +
		// 	"(?,?,?) "

		// _, err = tx.Exec(stmt2, id, id, 0)
		// if err != nil {
		// 	log.Panic(err)
		// }
		log.Println("Root Catalog ID: " + id)

	} else {

		stmt1 := "INSERT INTO catalog (id, version, name, locking, sequence, parent_id, created_by, created_time) " +
			"SELECT '" + id + "' as id, 0 as version, '" + child + "' as name, 0 as locking, 1 as sequence, " +
			"id as parent_id, 'admin' as created_by, NOW() as created_time " +
			"FROM catalog WHERE name='" + parents + "'"

		_, err = tx.Exec(stmt1)
		if err != nil {
			log.Panic(err)
		}

		// stmt2 := "INSERT INTO catalog_tree " +
		// 	"SELECT id as ancestor, '" + id + "' as descendant, " + depth +
		// 	" as depth FROM catalog WHERE name='" + parents + "'"

		// _, err = tx.Exec(stmt2)
		// if err != nil {
		// 	log.Panic(err)
		// }

		// stmtd := "INSERT INTO catalog_tree (ancestor, descendant, depth) VALUES (?,?,?) "
		// _, err = tx.Exec(stmtd, id, id, 0)
		// if err != nil {
		// 	log.Panic(err)
		// }
	}

	stmt0 := "INSERT INTO authority (id, item_id, item_type) VALUES (?,?,?)"
	_, err = tx.Exec(stmt0, "authority_"+id, id, "CATALOG")
	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	log.Println("The location was created.")
	return err
}

func CreateHierarchy(root string) error {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	log.Println("URL: " + url)
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var catalogList = []Catalog{}

	var no, id, name string
	rows, err := db.Query("WITH RECURSIVE tree as ( " +
		"SELECT id, name FROM catalog WHERE id='" + root + "' " +
		"UNION " +
		"SELECT catalog.id, catalog.name FROM catalog, tree WHERE tree.id=catalog.parent_id " +
		") " +
		"SELECT @rownum:=@rownum+1 No, t.* FROM tree t, (SELECT @rownum:=0) r ")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&no, &id, &name)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println("ID: " + id + ", Name: " + name)
		catalogList = append(catalogList, Catalog{no, id, name})
	}

	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	stmt2 := "INSERT INTO catalog_tree (ancestor, descendant, depth) VALUES (?,?,?) "

	for _, c1 := range catalogList {
		no, _ := strconv.Atoi(c1.No)
		idx := 0
		for i := no - 1; i < len(catalogList); i++ {

			if c1.Id == catalogList[i].Id {
				_, err = tx.Exec(stmt2, catalogList[i].Id, catalogList[i].Id, idx)
				if err != nil {
					log.Panic(err)
				}
			} else {
				_, err = tx.Exec(stmt2, c1.Id, catalogList[i].Id, idx)
				if err != nil {
					log.Panic(err)
				}
			}
			idx++
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	log.Println("The hierarchy was created.")
	return err
}

type Catalog struct {
	No   string
	Id   string
	Name string
}

func GetRoot() ([]Catalog, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var catalogList = []Catalog{}

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
		// log.Println("ID: " + id + ", Name: " + name)
		catalogList = append(catalogList, Catalog{"0", id, name})
	}

	return catalogList, err
}

func GetNodes(root string) ([]Catalog, error) {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var catalogList = []Catalog{}

	var id, name string
	rows, err := db.Query("WITH RECURSIVE tree as ( " +
		"SELECT id, name FROM catalog WHERE id='" + root + "' " +
		"UNION " +
		"SELECT catalog.id, catalog.name FROM catalog, tree WHERE tree.id=catalog.parent_id " +
		") " +
		"SELECT * FROM tree")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println("ID: " + id + ", Name: " + name)
		catalogList = append(catalogList, Catalog{"0", id, name})
	}

	return catalogList, err
}
