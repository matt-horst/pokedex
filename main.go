package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/matt-horst/pokeapi"
	"github.com/matt-horst/pokecache"
)

type cliCommand struct {
	name string
	description string
	callback func(*config) error
}

type config struct {
	Next string
	Previous string
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)

	return strings.Fields(text)
}

func commandHelp(_ *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print(usage)
	return nil
}

func commandExit(_ *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandMap(config *config) error {
	val, ok := cache.Get(config.Next)
	if ok {
		fmt.Print(string(val))
		return nil
	}

	list, err := pokeapi.GetLocationsList(config.Next)
	if err != nil {
		return err
	}

	output := ""
	for _, name := range list.Locations {
		output += name + "\n"
	}

	cache.Add(config.Next, []byte(output))
	fmt.Print(output)

	config.Next = list.Next
	config.Previous = list.Previous

	return nil
}

func commandMapb(config *config) error {
	val, ok := cache.Get(config.Previous)
	if ok {
		fmt.Print(string(val))
		return nil
	}

	list, err := pokeapi.GetLocationsList(config.Previous)
	if err != nil {
		return err
	}

	output := ""
	for _, name := range list.Locations {
		output += name + "\n"
	}

	cache.Add(config.Previous, []byte(output))
	fmt.Print(output)

	config.Next = list.Next
	config.Previous = list.Previous

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
	"exit": {
		name: "exit",
		description: "Exit the Pokedex",
		callback: commandExit,
	},
}

var cache = pokecache.NewCache(5 * time.Second)

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
			err := cmd.callback(&config)
			if err != nil {
				fmt.Println(err)
			} 
		} else {
			fmt.Println("Unknown command")
		}
	}
}
