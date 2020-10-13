package interactive

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

func PromptSensor(assetID string) (string, error) {

	validate := func(input string) error {
		url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
		db, err := sql.Open("mysql", url)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var id, name string
		rows, err := db.Query("SELECT id, name FROM sensor WHERE asset_id=? and name=?", assetID, input)

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
			return errors.New("duplicated sensor")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Sensor Name",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptSensorDuration() (string, error) {

	validate := func(input string) error {
		_, err := strconv.ParseUint(input, 0, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Sensor Duration(Sec.)",
		Validate: validate,
	}

	result, err := prompt.Run()
	r := "PT" + result + "S"
	return r, err
}
