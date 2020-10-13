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

func PromptParameter(paramGroupID string) (string, string, string, string, error) {

	validate := func(input string) error {
		r := strings.Split(input, ",")
		if len(r) < 3 {
			return nil
		}

		url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
		db, err := sql.Open("mysql", url)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var id, name string
		rows, err := db.Query("SELECT id, name FROM parameter WHERE param_group_id=? and name=?", paramGroupID, r[0])

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
			return errors.New("Duplicate parameter exists.")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Input a parameter [name,upper,target,lower]",
		Validate: validate,
	}

	result, err := prompt.Run()
	r := strings.Split(result, ",")
	return r[0], r[1], r[2], r[3], err
}

func PromptAddMoreParameter() (string, error) {
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
		Label:    "Do you want to add more parameters? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptAddParameter() (string, error) {
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
		Label:    "Would you like to add a parameter? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}
