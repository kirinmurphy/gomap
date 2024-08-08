package main

import (
	"gomap/src/locationHelpers"
	"log"
	"sync"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

type LocationStore struct {
	locations []locationHelpers.Location
	mu        sync.RWMutex
}

func (ls *LocationStore) GetLocations() []locationHelpers.Location {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	return ls.locations
}

func (ls *LocationStore) SetLocations(newLocations []locationHelpers.Location) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.locations = newLocations
}

func main() {
	initialLocations, err := locationHelpers.LoadLocations()
	if err != nil {
		log.Fatalf("failed to load initial locations: %v", err)
	}

	locationStore := &LocationStore{
		locations: initialLocations,
	}

	setLocations := func(newLocations []locationHelpers.Location) {
		locationStore.SetLocations(newLocations)
	}

	initRouter(locationStore, setLocations)
}
