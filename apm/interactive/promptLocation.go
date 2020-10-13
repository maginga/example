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

func PromptAddLocation() (string, error) {
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
		Label:    "Would you like to add a location? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptAddMoreLocation() (string, error) {
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
		Label:    "Do you want to add more locations? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptLocation() (string, string, string, error) {

	validate := func(input string) error {
		res := strings.Split(input, ",")
		if len(res) < 2 {
			return nil
		}

		url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
		db, err := sql.Open("mysql", url)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var id, name string
		rows, err := db.Query("SELECT id, name FROM catalog WHERE name=?", res[0])
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

		if len(id) > 0 {
			return errors.New("child name already exists.")
		}

		return err
	}

	prompt := promptui.Prompt{
		Label:    "Catalog Name [child, parents, depth]",
		Validate: validate,
	}

	result, err := prompt.Run()
	res := strings.Split(result, ",")
	return res[0], res[1], res[2], err
}
