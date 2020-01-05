package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	// TODO: Use modules instead of go get
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
		if !contains(InvalidPokemon, pokemon.Name) {
			validPokemonById[dexNumber] = pokemon
		}
	}

	sortedKeys := getSortedKeys(validPokemonById)
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

	// TODO: Move to a template
	fmt.Fprintf(writer, fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8">
			<meta http-equiv="x-ua-compatible" content="ie=edge">
			<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
			<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
			<title>Live shiny rates for Pokémon Go</title>

			<style>
				.icon {
					width: 60px;
					height: 60px;
				}
					
				#header {
					font-weight: bold;
					font-size: 25px;
					text-align: center;
					margin: 10px 0 0 0;
				}
				
				#data_period {
					font-size: 15px;
					text-align: center;
					margin: 10px;
				}
				
				.table > thead > tr > th {
					vertical-align: middle;
				}
			
				.table > tbody > tr > td {
					vertical-align: middle;
				}
				
				#footer {
					font-size: 13px;
					text-align: center;
					margin: 10px;
				}
			</style>
		</head>
		<body>
			<div id="header">
				Live Shiny Rates for Pokémon Go
			</div>
			<div id="data_period">
				<!-- Data from the last 24 hours. -->
			</div>
			<div id="shiny_table">
				<table class="table table-striped table-hover table-sm">
				<thead class="thead-dark">
					<tr>
					<th scope="col"/>
					<th scope="col">ID</th>
					<th scope="col">Name</th>
					<th scope="col">Shiny Rate</th>
					<th scope="col">Sample Size</th>
					</tr>
				</thead>
				<tbody id="table_body">
						%s
				</tbody>
				</table>
			</div>
			<div id="footer">
				Data is kindly provided by <a href="https://shinyrates.com">shinyrates.com</a>, updated every 6 hours.
			</div>
		</body>
	<html>
	`, shinyTableHtml))
}

func getSortedKeys(pokemonById map[int]PokemonDb) []int {
	keys := make([]int, 0, len(pokemonById))

	for key := range pokemonById {
		keys = append(keys, key)
	}

	sort.Ints(keys)
	return keys
}

func contains(array []string, input string) bool {
    for _, element := range array {
        if input == element {
            return true
        }
    }
    return false
}