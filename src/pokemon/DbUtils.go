package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

func OpenPokemonDb() (*sql.DB, error) {
	return openDb("postgres", "shiny_user", "shroot", "postgres", "shiny_db")
}

// TODO: If entry exists for day, UPDATE, else INSERT new entry
func DumpPokemon(db *sql.DB, pokemon Pokemon) error {
	_, err := db.Exec("INSERT INTO shiny_pokemon(pokemon, id, seen, found) VALUES($1, $2, $3, $4)", pokemon.Name, pokemon.DexNumber, pokemon.NumberSeen, pokemon.TotalFound)
	return err
}

func GetPokemonFromDb(db *sql.DB) (*sql.Rows, error) {
	return dbSelect(db, "shiny_pokemon", "pokemon", "id", "seen", "found")
}

// TODO: Extract to more generic DB class
// Stop hardcoding password
func openDb(driverType, dbUser, dbPassword, dbHost, dbName string) (*sql.DB, error) {
	return sql.Open(driverType, fmt.Sprintf("%1s://%2s:%3s@%4s/%5s?sslmode=disable", driverType, dbUser, dbPassword, dbHost, dbName))
}

func dbSelect(db *sql.DB, tableName string, columns ...string) (*sql.Rows, error) {
	columnString := strings.Join(columns[:], ",")
	return db.Query(fmt.Sprintf("SELECT %1s FROM %2s", columnString, tableName))
}
