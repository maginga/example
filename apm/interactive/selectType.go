package interactive

import (
	"example/apm/domain"

	"github.com/manifoldco/promptui"
)

func SelectType(tenantId string, roleName string) (string, string, error) {

	list, _ := domain.GetTypeList(tenantId, roleName)

	prompt := promptui.Select{
		Label: "Select Type",
		Items: list,
	}

	idx, _, err := prompt.Run()
	tuple := list[idx]
	return tuple.Id, tuple.Name, err
}
