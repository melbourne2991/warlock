package cmd

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"syscall"
)

func passwordPrompt() string {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	return string(bytePassword)
}
