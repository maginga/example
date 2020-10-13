package interactive

import (
	"example/apm/domain"

	"github.com/manifoldco/promptui"
)

func SelectParamGroup(tenantId string) (string, string, error) {

	list, _ := domain.GetParamGroup(tenantId)

	prompt := promptui.Select{
		Label: "Select Asset Type",
		Items: list,
	}

	idx, _, err := prompt.Run()
	tuple := list[idx]
	return tuple.Id, tuple.Name, err
}
