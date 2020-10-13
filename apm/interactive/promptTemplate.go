package interactive

import (
	"errors"

	"github.com/manifoldco/promptui"
)

func PromptTemplate() (string, error) {

	validate := func(input string) error {
		// _, err := strconv.ParseFloat(input, 64)
		// if err != nil {
		// 	return errors.New("Invalid number")
		// }
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Template Name",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptAddTemplate() (string, error) {
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
		Label:    "Would you like to add a template? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}
