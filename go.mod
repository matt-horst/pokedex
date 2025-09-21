module github.com/matt-horst/pokedex

go 1.25.1

require github.com/matt-horst/pokeapi v1.0.0
replace github.com/matt-horst/pokeapi => ./internal/pokeapi

require github.com/matt-horst/pokecache v1.0.0
replace github.com/matt-horst/pokecache => ./internal/pokecache
