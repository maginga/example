package interactive

import (
	"errors"

	"github.com/manifoldco/promptui"
)

func PromptTenant() (string, error) {

	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Tenant Name",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptAddTenant() (string, error) {
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
		Label:    "Would you like to add a tenant? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}
