package interactive

import (
	"example/apm/domain"

	"github.com/manifoldco/promptui"
)

func SelectOrg() (string, string, error) {

	list, _ := domain.GetOrgList()

	prompt := promptui.Select{
		Label: "Select Organization",
		Items: list,
	}

	idx, _, err := prompt.Run()
	tuple := list[idx]
	return tuple.Id, tuple.Name, err
}
