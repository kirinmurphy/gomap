package router

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func InitRouter(redisClient *redis.Client, ctx context.Context) {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeRouteHandler(w, r, redisClient, ctx)
	})

	r.HandleFunc("/loadLocations", func(w http.ResponseWriter, r *http.Request) {
		loadLocationsRouteHandler(w, r, redisClient, ctx)
	}).Methods("GET")

	r.HandleFunc("/updateMapUI", func(w http.ResponseWriter, r *http.Request) {
		updateMapUIHandler(w, r, redisClient, ctx)
	}).Methods("POST")

	r.HandleFunc("/getLocations", func(w http.ResponseWriter, r *http.Request) {
		getLocationsRouteHandler(w, r, redisClient)
	})

	r.PathPrefix("/src/templates/").Handler(
		http.StripPrefix("/src/templates/", http.FileServer(http.Dir("src/templates"))),
	)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
