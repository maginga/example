package interactive

import (
	"errors"

	"github.com/manifoldco/promptui"
)

func PromptAddDevice() (string, error) {
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
		Label:    "Would you like to add a device? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptAddMoreDevice() (string, error) {
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
		Label:    "Do you want to add more device? (y/n)",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptDeviceAddress() (string, error) {

	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Device's IP address",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptDeviceMacAddress() (string, error) {

	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Device's MAC address",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptDeviceName() (string, error) {

	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Device Name",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptDeviceModelNumber() (string, error) {

	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Device's Model Number",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}

func PromptDeviceSerialNumber() (string, error) {

	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Device's Serial Number",
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
}
