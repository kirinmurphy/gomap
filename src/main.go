package main

import (
	"context"
	"gomap/src/locationHelpers"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

func init() {
	godotenv.Load()

	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	log.Println("Connected to redis")
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
	// initialLocations, err := locationHelpers.LoadLocations()
	// if err != nil {
	// 	log.Fatalf("failed to load initial locations: %v", err)
	// }

	// locationStore := &LocationStore{
	// 	locations: initialLocations,
	// }

	// setLocations := func(newLocations []locationHelpers.Location) {
	// 	locationStore.SetLocations(newLocations)
	// }

	// locationStore := &LocationStore{}
	initRouter()
}
