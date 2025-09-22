package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"errors"
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
		return errors.New("failed to retrieve map information")
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
		return errors.New("failed to retrieve map information")
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
		return fmt.Errorf("failed to retrieve location information for `%v`", name)
	}

	fmt.Printf("Exploring %v...\n", name)
	for _, p := range pokemon {
		fmt.Printf(" - %v\n", p)
	}

	return nil
}

func commandCatch(_ *config, params []string) error {
	name := params[0]

	p, err := pokeapi.GetPokemon(name)
	if err != nil {
		return fmt.Errorf("failed to retrieve pokemon information for `%v`", name)
	}


	fmt.Printf("Throwing a Pokeball at %v...\n", name)

	rng := rand.New(rand.NewSource(int64(time.Now().UnixNano())))
	v := rng.Intn(400)
	if v > p.BaseExperience {
		// Catch
		fmt.Printf("%v was caught!\n", name)
		pokedex[name] = p
	} else {
		// Miss
		fmt.Printf("%v escaped!\n", name)
	}

	return nil
}

func commandInspect(_ *config, params []string) error {
	name := params[0]

	p, ok := pokedex[name]
	if !ok {
		fmt.Printf("You must catch a %v before inspecting\n", name)
		return nil
	}

	fmt.Printf("Name: %v\n", p.Name)
	fmt.Printf("Height: %v\n", p.Height)
	fmt.Printf("Weight: %v\n", p.Weight)

	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf(" -%v: %v\n", s.Name, s.Val)
	}

	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf(" - %v\n", t)
	}

	return nil
}

func commandPokedex(_ *config, _ []string) error {
	fmt.Println("Your Pokedex:")

	for p := range pokedex {
		fmt.Printf(" - %v\n", p)
	}

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
	"catch": {
		name: "catch",
		description: "Attempts to catch a pokemon",
		callback: commandCatch,
	},
	"inspect": {
		name: "inspect",
		description: "Displays information about pokemon in pokedex",
		callback: commandInspect,
	},
	"pokedex": {
		name: "pokedex",
		description: "Lists all known pokemon in your pokedex",
		callback: commandPokedex,
	},
	"exit": {
		name: "exit",
		description: "Exit the Pokedex",
		callback: commandExit,
	},
}

var pokedex map[string]pokeapi.PokemonInfo = make(map[string]pokeapi.PokemonInfo)

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
				fmt.Printf("Error: %v\n", err)
			} 
		} else {
			fmt.Println("Unknown command")
		}
	}
}
