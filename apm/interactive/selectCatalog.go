package interactive

import (
	"example/apm/domain"

	"github.com/manifoldco/promptui"
)

func SelectRootCatalog() (string, string, error) {

	catalogs, _ := domain.GetRoot()

	prompt := promptui.Select{
		Label: "Select Root Catalog",
		Items: catalogs,
	}

	idx, _, err := prompt.Run()
	catalog := catalogs[idx]
	return catalog.Id, catalog.Name, err
}

func SelectCatalogs(root string) (string, string, error) {

	catalogs, _ := domain.GetNodes(root)

	prompt := promptui.Select{
		Label: "Select Catalog",
		Items: catalogs,
	}

	idx, _, err := prompt.Run()
	catalog := catalogs[idx]
	return catalog.Id, catalog.Name, err
}
