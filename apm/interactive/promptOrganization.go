package interactive

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

// var logger *log.Logger

func PromptOrganization() (string, error) {
	// logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)

	validate := func(input string) error {
		url := fmt.Sprintf("%v", viper.Get("metadata.discovery-url"))
		db, err := sql.Open("mysql", url)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var id string
		rows, err := db.Query("SELECT id FROM user_org WHERE org_name=?", strings.ToUpper(input)+"_ORG")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&id)
			if err != nil {
				log.Fatal(err)
			}
		}
		if len(id) > 0 {
			return errors.New("Organization name already exists.")
		}

		return err
	}

	prompt := promptui.Prompt{
		Label:    "Organization Name",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptAddOrganization() (string, error) {
	validate := func(input string) error {
		var err error
		if input == "y" || input == "n" {
			err = nil
		} else {
			err = errors.New("Invalid input")
		}

		return err
	}

	prompt := promptui.Prompt{
		Label:    "Would you like to add a organization? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}
