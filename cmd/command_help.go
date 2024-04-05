package cmd

import "fmt"

func helpCmd(cfg *config) error {
	fmt.Println("\nWelcome to Gymplanr")
	fmt.Println("Here are the available commands:")
	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf(" - %s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println("")

	return nil
}
