package interactive

import (
	"errors"

	"github.com/manifoldco/promptui"
)

func PromptParamGroup() (string, error) {

	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Parameter Group Name",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptAddParamGroup() (string, error) {
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
		Label:    "Would you like to add a parameter group? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}
