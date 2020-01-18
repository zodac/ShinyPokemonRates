package main

import (
	"sort"
)

func GetSortedKeys(pokemonById map[int]PokemonDb) []int {
	keys := make([]int, 0, len(pokemonById))

	for key := range pokemonById {
		keys = append(keys, key)
	}

	sort.Ints(keys)
	return keys
}

func Contains(array []string, input string) bool {
    for _, element := range array {
        if input == element {
            return true
        }
    }
    return false
}