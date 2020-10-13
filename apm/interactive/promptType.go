package interactive

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

func PromptType(tenantID string, roleName string) (string, error) {

	validate := func(input string) error {
		url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
		db, err := sql.Open("mysql", url)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var id, name string
		rows, err := db.Query("SELECT id, name FROM type WHERE tenant_id=? and role=? and name=?",
			tenantID, roleName, input)

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&id, &name)
			if err != nil {
				log.Fatal(err)
			}
		}

		if len(id) <= 0 {
			return nil
		}

		return errors.New("Duplicate value exists.")
	}

	prompt := promptui.Prompt{
		Label:    "Type Name [" + roleName + "]",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptAddMoreType() (string, error) {
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
		Label:    "Do you want to add more type? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptAddType(title string) (string, error) {
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
		Label:    "Would you like to add a " + title + " type? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}
