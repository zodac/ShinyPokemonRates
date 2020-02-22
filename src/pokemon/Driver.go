package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	humanize "github.com/dustin/go-humanize"
)

const (
	port             = "5000"
	templateRootDir  = "html/"
	templateFileName = "template.html"
	templateFilePath = templateRootDir + templateFileName
)

var (
	intervalInHours = GetEnv("INTERVAL_IN_HOURS", "24")
)

func main() {
	http.HandleFunc("/", showRates)
	http.HandleFunc("/manual", CheckPokemonManual)
	http.Handle("/html", http.FileServer(http.Dir(templateRootDir)))

	go CreateCronJob()
	fmt.Printf("Starting HTTP server on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// TODO: Split this into more classes
func showRates(writer http.ResponseWriter, request *http.Request) {
	db, err := OpenPokemonDb()
	if err != nil {
		fmt.Printf("Error opening DB connection: %s", err)
		return
	}

	rows, err := GetPokemonFromDb(db)
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
			ID:    id,
			Seen:  totalSeen,
			Found: totalFound,
		}

		pokemonFromDb[id] = newPokemon
	}

	validPokemonByID := make(map[int]PokemonDb)
	for dexNumber, pokemon := range pokemonFromDb {
		if !Contains(InvalidPokemon, pokemon.Name) {
			validPokemonByID[dexNumber] = pokemon
		}
	}

	sortedKeys := GetSortedKeys(validPokemonByID)
	shinyTableHTML := ""

	for _, sortedKey := range sortedKeys {
		sortedPokemon := validPokemonByID[sortedKey]
		rate := sortedPokemon.Seen / sortedPokemon.Found

		shinyTableHTML += fmt.Sprintf(`
						<tr>
							<th scope="row" id="%[1]s">
								<img class="icon" src="res/sprites/%[1]s.png">
							</th>
							<td>%[1]s</td>
							<td>%[2]s</td>
							<td>1/%[3]s</td>
							<td>%[4]s in %[5]s</td>
						</tr>
		`,
			strconv.Itoa(sortedPokemon.ID),
			sortedPokemon.Name,
			humanize.Comma(int64(rate)),
			humanize.Comma(int64(sortedPokemon.Found)),
			humanize.Comma(int64(sortedPokemon.Seen)))
	}

	// TODO: Add dividers for different gens
	// TODO: Use a range in the template?

	pageData := PageData{
		TableBody:       shinyTableHTML,
		IntervalInHours: intervalInHours,
	}

	templates := template.Must(template.ParseFiles(templateFilePath))
	templates.ExecuteTemplate(writer, templateFileName, pageData)
}

type PageData struct {
	TableBody       string
	IntervalInHours string
}
