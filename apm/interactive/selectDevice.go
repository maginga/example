package interactive

import (
	"example/apm/domain"

	"github.com/manifoldco/promptui"
)

func SelectDevice(tenantId string) (string, string, error) {

	list, _ := domain.GetDeviceList(tenantId)

	prompt := promptui.Select{
		Label: "Select Device",
		Items: list,
	}

	idx, _, err := prompt.Run()
	tuple := list[idx]
	return tuple.Id, tuple.Name, err
}
