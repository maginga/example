package interactive

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

func PromptAddMoreAsset() (string, error) {
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
		Label:    "Do you want to add more asset? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err

}

func PromptAsset() (string, error) {

	validate := func(input string) error {
		url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
		db, err := sql.Open("mysql", url)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var id, name string
		rows, err := db.Query("SELECT id, name FROM asset WHERE name=?", input)

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
			return errors.New("Duplicate asset exists.")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Input a asset",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}
