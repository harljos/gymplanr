package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name         string
	description string
	callback     func(*config) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    helpCmd,
		},
		"generate": {
			name:        "generate",
			description: "Generates workout plan",
			callback:    generateCmd,
		},
		"exit": {
			name:         "exit",
			description: "Exit gymplanr",
			callback:     exitCmd,
		},
	}
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("gymplanr >")

		scanner.Scan()
		text := scanner.Text()

		cleaned := cleanInput(text)
		if len(cleaned) == 0 {
			continue
		}

		commandName := cleaned[0]
		
		commands := getCommands()

		command, ok := commands[commandName]
		if !ok {
			fmt.Println("invalid command")
			continue
		}

		err := command.callback(cfg)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(str string) []string {
	lowered := strings.ToLower(str)
	words := strings.Fields(lowered)
	return words
}
