package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"gomap/src/router"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

const BASE_SPREADSHEET_URL = "https://docs.google.com/spreadsheets/d/e/%s/pub?gid=0&single=true&output=csv"

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

func main() {
	r := router.InitRouter(router.RouterConfig{
		RedisClient:        redisClient,
		Ctx:                ctx,
		BaseSpreadsheetUrl: BASE_SPREADSHEET_URL,
	})

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
