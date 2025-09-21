package pokeapi

import (
	"net/http"
	"io"
	"encoding/json"
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
