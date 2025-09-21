package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/matt-horst/pokeapi"
)

type cliCommand struct {
	name string
	description string
	callback func(*config, []string) error
}

type config struct {
	Next string
	Previous string
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)

	return strings.Fields(text)
}

func commandHelp(_ *config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print(usage)
	return nil
}

func commandExit(_ *config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandMap(config *config, _ []string) error {
	list, err := pokeapi.GetLocationsList(config.Next)
	if err != nil {
		return err
	}

	for _, name := range list.Locations {
		fmt.Println(name)
	}

	config.Next = list.Next
	config.Previous = list.Previous

	return nil
}

func commandMapb(config *config, _ []string) error {
	list, err := pokeapi.GetLocationsList(config.Previous)
	if err != nil {
		return err
	}

	for _, name := range list.Locations {
		fmt.Println(name)
	}

	config.Next = list.Next
	config.Previous = list.Previous

	return nil
}

func commandExplore(_ *config, params []string) error {
	name := params[0]

	pokemon, err := pokeapi.GetPokemonList(name)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %v...\n", name)
	for _, p := range pokemon {
		fmt.Printf(" - %v\n", p)
	}

	return nil
}

func commandCatch(_ *config, params []string) error {
	name := params[0]

	_, err := pokeapi.GetPokemon(name)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", name)

	return nil
}

var registry = map[string]cliCommand{
	"help": {
		name: "help",
		description: "Displays a help message",
		callback: commandHelp,
	},
	"map": {
		name: "map",
		description: "Displays a list of the next 20 locations",
		callback: commandMap,
	},
	"mapb": {
		name: "mapb",
		description: "Displays a list of the previous 20 locations",
		callback: commandMapb,
	},
	"explore": {
		name: "explore",
		description: "Displays a list of pokemon at the given location",
		callback: commandExplore,
	},
	"exit": {
		name: "exit",
		description: "Exit the Pokedex",
		callback: commandExit,
	},
}

type Pokemon struct {
	name string
}
var pokedex map[string]Pokemon = make(map[string]Pokemon)

var usage string
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
	config := config{}

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		rawInput := scanner.Text()
		cleanInput := cleanInput(rawInput)

		if cmd, ok := registry[cleanInput[0]]; ok {
			err := cmd.callback(&config, cleanInput[1:])
			if err != nil {
				fmt.Println(err)
			} 
		} else {
			fmt.Println("Unknown command")
		}
	}
}
