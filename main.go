package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)

	return strings.Fields(text)
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

var registry = map[string]cliCommand{
	"help": {
		name: "help",
		description: "Displays a help message",
		callback: commandHelp,
	},
	"exit": {
		name: "exit",
		description: "Exit the Pokedex",
		callback: commandExit,
	},
}

var usage string

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print(usage)
	return nil
}

func generateUsage() string {
	usage := "Usage:\n\n"
	for _, cmd := range registry {
		usage += fmt.Sprintf("%s: %s\n", cmd.name, cmd.description)
	}

	return usage
}


func main() {
	usage = generateUsage()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		rawInput := scanner.Text()
		cleanInput := cleanInput(rawInput)

		if cmd, ok := registry[cleanInput[0]]; ok {
			err := cmd.callback()
			if err != nil {
				fmt.Println(err)
			} 
		} else {
			fmt.Println("Unknown command")
		}
	}
}
