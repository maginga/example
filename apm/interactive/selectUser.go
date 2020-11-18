package interactive

import (
	"errors"
	"example/apm/domain"

	"github.com/manifoldco/promptui"
)

func PromptAddUser() (string, error) {
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
		Label:    "Would you like to add a user? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func SelectUser() (string, string, error) {

	list, _ := domain.GetUserList()

	prompt := promptui.Select{
		Label: "Select User",
		Items: list,
	}

	idx, _, err := prompt.Run()
	tuple := list[idx]
	return tuple.Id, tuple.Name, err
}
