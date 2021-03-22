package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	_ "github.com/denisenkom/go-mssqldb"
	"gopkg.in/yaml.v2"
)

func main() {
	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	log.Printf("MSSQL Conn: %s\n", config.ConnString)
	db, err := sql.Open("sqlserver", config.ConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var val int
	err = db.QueryRow("SELECT count(*) FROM dbo.history").Scan(&val)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)

	// var val int
	// rows, err := db.Query("SELECT count(*) FROM dbo.history")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	err := rows.Scan(&val)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(val)
	// }
}
