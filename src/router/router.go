package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type RedisClientInterface interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

var templateDir = "src/templates"

func InitRouter(redisClient RedisClientInterface, ctx context.Context) *mux.Router {
	r := mux.NewRouter()

	InitializeHomePageTemplates(templateDir)
	InitializeUpdateMapUITemplates(templateDir)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeRouteHandler(w, r, redisClient, ctx)
	})

	r.HandleFunc("/updateMapUI", func(w http.ResponseWriter, r *http.Request) {
		updateMapUIHandler(w, r, redisClient, ctx)
	}).Methods("POST")

	r.HandleFunc("/loadLocations", func(w http.ResponseWriter, r *http.Request) {
		loadLocationsRouteHandler(w, r, redisClient, ctx)
	}).Methods("GET")

	r.HandleFunc("/getLocations", func(w http.ResponseWriter, r *http.Request) {
		getLocationsRouteHandler(w, r, redisClient)
	})

	r.PathPrefix("/src/templates/").Handler(
		http.StripPrefix("/src/templates/", http.FileServer(http.Dir("src/templates"))),
	)

	return r
}
