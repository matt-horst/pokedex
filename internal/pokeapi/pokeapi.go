package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationsList struct {
	Locations []string
	Next string
	Previous string
}

func GetLocationsList(url string) (LocationsList, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	res, err := http.Get(url)

	if err != nil {
		return LocationsList{}, err
	} 

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return LocationsList{}, err
	}

	data := struct {
		Count int
		Next string
		Previous string
		Results []struct {
			Name string
			Url string
		}
	}{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return LocationsList{}, err
	}

	lst := LocationsList{Next: data.Next, Previous: data.Previous}

	for _, r := range data.Results {
		lst.Locations = append(lst.Locations, r.Name)
	}

	return lst, nil
}

func GetPokemonList(name string) ([]string, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", name)

	res, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return []string{}, err
	}

	data := struct {
		Pokemon_encounters []struct {
			Pokemon struct {
				Name string
			}
		}
	} {}

	if err := json.Unmarshal(body, &data); err != nil {
		return []string{}, err
	}

	pokemon := []string{}
	for _, p := range data.Pokemon_encounters {
		pokemon = append(pokemon, p.Pokemon.Name)
	}

	return pokemon, nil
}
