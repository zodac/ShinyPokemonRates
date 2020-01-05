package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"html/template"

	humanize "github.com/dustin/go-humanize"
)

const (
	port = "5000"
)

func main() {
	go CreateCronJob()
	fmt.Printf("Starting HTTP server on port %s\n", port)
	http.HandleFunc("/shiny", showRates)
	http.ListenAndServe(":"+port, nil)
}

// TODO: Split this into more classes
func showRates(writer http.ResponseWriter, request *http.Request) {
	db, err := sql.Open("postgres", "postgres://shiny_user:shroot@localhost/shiny_db?sslmode=disable")
	if err != nil {
		fmt.Printf("Error opening DB connection: %s", err)
		return
	}

	rows, err := db.Query("SELECT pokemon, id, seen, found FROM shiny_pokemon")
	if err != nil {
		fmt.Printf("Error reading from DB: %s\n", err)
		return
	}
	defer rows.Close()

	pokemonFromDb := make(map[int]PokemonDb)

	for rows.Next() {
		var pokemon string
		var id int
		var seen int
		var found int

		err = rows.Scan(&pokemon, &id, &seen, &found)
		if err != nil {
			fmt.Printf("Error scanning from DB: %s\n", err)
			continue
		}

		totalFound := found
		totalSeen := seen
		if existingPokemon, ok := pokemonFromDb[id]; ok {
			totalSeen += existingPokemon.Seen
			totalFound += existingPokemon.Found
		}

		newPokemon := PokemonDb{
			Name:  pokemon,
			Id:    id,
			Seen:  totalSeen,
			Found: totalFound,
		}

		pokemonFromDb[id] = newPokemon
	}

	validPokemonById := make(map[int]PokemonDb)
	for dexNumber, pokemon := range pokemonFromDb {
		if !Contains(InvalidPokemon, pokemon.Name) {
			validPokemonById[dexNumber] = pokemon
		}
	}

	sortedKeys := GetSortedKeys(validPokemonById)
	shinyTableHtml := ""

	for _, sortedKey := range sortedKeys {
		sortedPokemon := validPokemonById[sortedKey]
		rate := sortedPokemon.Seen/sortedPokemon.Found

		// TODO: Use <img class="icon" src="./icons/%[1]s.png"> when possible
		shinyTableHtml += fmt.Sprintf(`
						<tr>
							<th scope="row" id="%[1]s">
								<img class="icon" src="https://shinyrates.com/images/icons/%[1]s.png">
							</th>
							<td>%[1]s</td>
							<td>%[2]s</td>
							<td>1/%[3]s</td>
							<td>%[4]s in %[5]s</td>
						</tr>
		`,
		strconv.Itoa(sortedPokemon.Id),
		sortedPokemon.Name,
		humanize.Comma(int64(rate)),
		humanize.Comma(int64(sortedPokemon.Found)),
		humanize.Comma(int64(sortedPokemon.Seen)))
	}

	// TODO: Add dividers for different gens
	template, _ := template.ParseFiles("html/index.html")
	template.Execute(writer, shinyTableHtml)
}