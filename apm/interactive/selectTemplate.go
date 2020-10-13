package interactive

import (
	"example/apm/domain"

	"github.com/manifoldco/promptui"
)

func SelectTemplate(tenantId string) (string, string, error) {

	list, _ := domain.GetTemplateList(tenantId)

	prompt := promptui.Select{
		Label: "Select Template",
		Items: list,
	}

	idx, _, err := prompt.Run()
	tuple := list[idx]
	return tuple.Id, tuple.Name, err
}
