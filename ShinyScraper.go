package main

import (
	"fmt"
	"strings"
	"time"

	"database/sql"
	_ "github.com/lib/pq"

	"github.com/jasonlvhit/gocron"
	"sourcegraph.com/sourcegraph/go-selenium"
)

const (
	url                             = "https://shinyrates.com"
	cronIntervalInHours = 6
	dockerHost = "localhost"
)

func CreateCronJob() {
	cronJob := gocron.NewScheduler()
	cronJob.Every(cronIntervalInHours).Hours().Do(checkPokemon)

	fmt.Printf("Checking Pokemon every %v hours...\n", cronIntervalInHours)
	checkPokemon() // Check once before the CRON job starts
	<-cronJob.Start()
}

func checkPokemon() {
	allPokemonById, err := getAllPokemonById()
	if err != nil {
		fmt.Printf("Failed to retrieve all Pokemon: %s", err)
		return
	}

	currentTime := time.Now()
	fmt.Printf("\nFound (%v) filtered Pokemon at %s\n", len(allPokemonById), currentTime.Format("2006-01-02 15:04"))

	db, err := sql.Open("postgres", "postgres://shiny_user:shroot@"+dockerHost+"/shiny_db?sslmode=disable")
	if err != nil {
		fmt.Printf("Error opening DB connection: %s", err)
		return
	}

	for _, pokemon := range allPokemonById {
		err := dumpPokemon(db, pokemon)
		if err != nil {
			fmt.Printf("Error %s: %s\n", pokemon.Name, err)
		}
	}
	fmt.Println("Dumped all Pokemon to DB")
}

func dumpPokemon(db *sql.DB, pokemon Pokemon) error {
	_, err := db.Exec("INSERT INTO shiny_pokemon(pokemon, id, seen, found) VALUES($1, $2, $3, $4)", pokemon.Name, pokemon.DexNumber, pokemon.NumberSeen, pokemon.TotalFound)
	return err
}

func getAllPokemonById() (map[int]Pokemon, error) {
	shinyTable, err := getShinyTable()
	if err != nil {
		fmt.Printf("Failed to retrieve shiny table: %s", err)
		return nil, err
	}

	shinyTableEntries := strings.Split(shinyTable, "\n")

	pokemonById := make(map[int]Pokemon)
	for _, shinyTableEntry := range shinyTableEntries[1:] {
		pokemon := NewPokemon(shinyTableEntry)
		pokemonById[pokemon.DexNumber] = pokemon
	}

	return pokemonById, nil
}

func getShinyTable() (string, error) {
	capabilities := selenium.Capabilities(map[string]interface{}{
		"browserName": "chrome",
	})

	webDriver, err := selenium.NewRemote(capabilities, "http://"+dockerHost+":4444/wd/hub")
	if err != nil {
		fmt.Printf("Failed to open session: %s", err)
		return "", err
	}
	defer webDriver.Quit()

	err = webDriver.Get(url)
	if err != nil {
		fmt.Printf("Failed to load page: %s", err)
		return "", err
	}

	time.Sleep(10 * time.Second)

	elem, err := webDriver.FindElement(selenium.ByCSSSelector, "#shiny_table")
	if err != nil {
		fmt.Printf("Failed to find elements: %s", err)
		return "", err
	}

	return elem.Text()
}
