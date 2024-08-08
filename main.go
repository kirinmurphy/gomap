package main

import (
	"log"

	"gomap/locationHelpers"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	locations, err := locationHelpers.LoadLocations()
	if err != nil {
		log.Fatalf("failed to load initial locations: %v", err)
	}

	initRouter(locations)
}
