package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jasonlvhit/gocron"
	"sourcegraph.com/sourcegraph/go-selenium"
)

const (
	url          = "https://shinyrates.com"
	dockerHost   = "localhost"
	seleniumHost = "selenium-hub"
)

var (
	cronIntervalInHoursString = GetEnv("INTERVAL_IN_HOURS", "24")
)

func CreateCronJob() {
	cronJob := gocron.NewScheduler()

	cronIntervalInHours, err := strconv.ParseInt(cronIntervalInHoursString, 10, 64)
	if err != nil {
		fmt.Printf("Cannot determine interval period: %s", err)
		panic(fmt.Sprintf("Invalid value for 'INTERVAL_IN_HOURS' provided, expected hours, found: '%s'", cronIntervalInHoursString))
	}

	// Check once before the CRON job starts, add delay to allow Selenium to come up
	time.Sleep(10 * time.Second)
	CheckPokemon()

	cronJob.Every(uint64(cronIntervalInHours)).Hours().Do(CheckPokemon)
	fmt.Printf("Checking Pokemon every %v hours...\n", cronIntervalInHours)
	<-cronJob.Start()
}

func CheckPokemonManual(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Checking for Pokemon")
	CheckPokemon()
	fmt.Fprintf(writer, "Pokemon loaded!")
}

func CheckPokemon() {
	allPokemonByID, err := getAllPokemonByID()
	if err != nil {
		fmt.Printf("Failed to retrieve all Pokemon: %s", err)
		return
	}

	currentTime := time.Now()
	fmt.Printf("\nFound (%v) filtered Pokemon at %s\n", len(allPokemonByID), currentTime.Format("2006-01-02 15:04"))

	db, err := OpenPokemonDb()
	if err != nil {
		fmt.Printf("Error opening DB connection: %s", err)
		return
	}

	for _, pokemon := range allPokemonByID {
		err := DumpPokemon(db, pokemon)
		if err != nil {
			fmt.Printf("Error %s: %s\n", pokemon.Name, err)
		}
	}
	fmt.Println("Dumped all Pokemon to DB")
}

func getAllPokemonByID() (map[int]Pokemon, error) {
	shinyTable, err := getShinyTable()
	if err != nil {
		fmt.Printf("Failed to retrieve shiny table: %s", err)
		return nil, err
	}

	shinyTableEntries := strings.Split(shinyTable, "\n")

	pokemonByID := make(map[int]Pokemon)
	for _, shinyTableEntry := range shinyTableEntries[1:] {
		pokemon := NewPokemon(shinyTableEntry)
		pokemonByID[pokemon.DexNumber] = pokemon
	}

	return pokemonByID, nil
}

func getShinyTable() (string, error) {
	capabilities := selenium.Capabilities(map[string]interface{}{
		"browserName": "chrome",
	})

	webDriver, err := selenium.NewRemote(capabilities, "http://"+seleniumHost+":4444/wd/hub")
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
