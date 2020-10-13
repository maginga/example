package interactive

import (
	"example/apm/domain"

	"github.com/manifoldco/promptui"
)

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
