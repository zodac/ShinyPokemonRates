package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Pokemon struct {
	Name       string
	DexNumber  int
	CatchRate  string
	NumberSeen string
	TotalFound string
}

type PokemonDb struct {
	Name  string
	ID    int
	Seen  int
	Found int
}

func (p *Pokemon) String() string {
	return fmt.Sprintf("%12s (%3v): \tFound %3s in %6s\t(1 in %5s)", p.Name, p.DexNumber, p.TotalFound, p.NumberSeen, p.CatchRate)
}

func (p *PokemonDb) String() string {
	return fmt.Sprintf("%12s:\tFound %3v in %6v", p.Name, p.Found, p.Seen)
}

//	1	Bulbasaur	1/474	3,319
func NewPokemon(shinyTableEntry string) Pokemon {
	pokemonElements := strings.Fields(shinyTableEntry)
	dexNumber, _ := strconv.Atoi(pokemonElements[0])

	catchRateTokens := strings.Split(pokemonElements[2], "/")
	catchRateRaw := strings.Replace(catchRateTokens[1], ",", "", -1)
	numberSeenRaw := strings.Replace(pokemonElements[3], ",", "", -1)

	catchRate, _ := strconv.Atoi(catchRateRaw)
	numberSeen, _ := strconv.Atoi(numberSeenRaw)
	totalFound := strconv.Itoa(numberSeen / catchRate)

	pokemon := Pokemon{
		Name:       pokemonElements[1],
		DexNumber:  dexNumber,
		CatchRate:  catchRateRaw,
		NumberSeen: numberSeenRaw,
		TotalFound: totalFound,
	}

	return pokemon
}

// TODO: Move to file
var InvalidPokemon = []string{
	// Gen 1
	"Bulbasaur",
	"Charmander",
	"Squirtle",
	"Caterpie",
	"Pidgey",
	"Pikachu",
	"Nidoran♂",
	"Oddish",
	"Diglett",
	"Psyduck",
	"Growlithe",
	"Poliwag",
	"Machop",
	"Tentacool",
	"Geodude",
	"Ponyta",
	"Magnemite",
	"Seel",
	"Shellder",
	"Gastly",
	"Onix",
	"Krabby",
	"Cubone",
	"Koffing",
	"Kangaskhan",
	"Horsea",
	"Mr. Mime",
	"Scyther",
	"Jynx",
	"Pinsir",
	"Tauros",
	"Magikarp",
	"Eevee",
	"Omanyte",
	"Kabuto",
	"Dratini",

	// Gen 2
	"Chikorita",
	"Cyndaquil",
	"Totodile",
	"Sentret",
	"Natu",
	"Mareep",
	"Aipom",
	"Yanma",
	"Sunkern",
	"Murkrow",
	"Misdreavus",
	"Pineco",
	"Gligar",
	"Snubbull",
	"Shuckle",
	"Sneasel",
	"Swinub",
	"Delibird",
	"Skarmory",
	"Houndour",
	"Larvitar",

	// Gen 3
	"Treecko",
	"Torchic",
	"Mudkip",
	"Poochyena",
	"Zigzagoon",
	"Lotad",
	"Wingull",
	"Ralts",
	"Slakoth",
	"Makuhita",
	"Sableye",
	"Mawile",
	"Aron",
	"Meditite",
	"Electrike",
	"Plusle",
	"Minun",
	"Roselia",
	"Carvanha",
	"Wailmer",
	"Spoink",
	"Trapinch",
	"Swablu",
	"Zangoose",
	"Seviper",
	"Solrock",
	"Barboach",
	"Lileep",
	"Anorith",
	"Feebas",
	"Castform",
	"Shuppet",
	"Duskull",
	"Snorunt",
	"Clamperl",
	"Luvdisc",
	"Bagon",
	"Beldum",

	// Gen 4
	"Turtwig",
	"Chimchar",
	"Piplup",
	"Shinx",
	"Drifloon",
	"Buneary",
	"Snover",

	// Gen 5
	"Patrat",
	"Lillipup",
}
