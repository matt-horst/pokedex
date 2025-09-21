package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/matt-horst/pokecache"
	"time"
)

type LocationsList struct {
	Locations []string
	Next string
	Previous string
}

var cache *pokecache.Cache = pokecache.NewCache(5 * time.Second)

func GetLocationsList(url string) (LocationsList, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	body, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)

		if err != nil {
			return LocationsList{}, err
		} 

		body, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return LocationsList{}, err
		}

		cache.Add(url, body)
	}

	data := struct {
		Next string
		Previous string
		Results []struct {
			Name string
		}
	}{}

	err := json.Unmarshal(body, &data)
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

	body, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return []string{}, err
		}

		body, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return []string{}, err
		}

		cache.Add(url, body)
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
