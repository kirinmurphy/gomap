package main

import (
	"context"
	"log"
	"os"

	"gomap/src/router"

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

func main() {
	router.InitRouter(redisClient, ctx)
}
