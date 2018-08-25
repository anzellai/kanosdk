package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

const (
	actionComm = "Comm"
	actionQuit = "Quit"
)

func devicePrompt() (string, error) {
	device, err := makePrompt("Please enter DEVICE address to connect to")
	if err != nil {
		return "", err
	}
	return device, nil
}

func commPrompt() (string, error) {
	data, err := makePrompt("Please enter DATA to send (q/quit = QUIT)")
	if err != nil {
		return "", err
	}
	return data, nil
}

func makePrompt(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}
	fmt.Printf("You enter %q\n", result)
	return result, nil
}
