package interactive

import (
	"example/apm/domain"

	"github.com/manifoldco/promptui"
)

func SelectNest(tenantId string) (string, string, error) {

	list, _ := domain.GetNestList(tenantId)

	prompt := promptui.Select{
		Label: "Select Nest",
		Items: list,
	}

	idx, _, err := prompt.Run()
	tuple := list[idx]
	return tuple.Id, tuple.Name, err
}

func SelectTenant() (string, string, error) {

	list, _ := domain.GetTenantList()

	prompt := promptui.Select{
		Label: "Select Tenant",
		Items: list,
	}

	idx, _, err := prompt.Run()
	tuple := list[idx]
	return tuple.Id, tuple.Name, err
}
