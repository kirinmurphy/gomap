package main

import (
	"encoding/json"
	"fmt"
	"gomap/src/locationHelpers"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

const baseSpreadsheetUrl = "https://docs.google.com/spreadsheets/d/e/%s/pub?gid=0&single=true&output=csv"

func initRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)

	r.HandleFunc("/loadLocations", loadLocationsHandler).Methods("GET")

	r.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
		locationsHandler(w, r)
	})

	// r.HandleFunc("/reload-locations", func(w http.ResponseWriter, r *http.Request) {
	// 	reloadLocationsHandler(w, r, setLocations, locationStore)
	// }).Methods("GET")

	r.PathPrefix("/src/templates/").Handler(
		http.StripPrefix("/src/templates/", http.FileServer(http.Dir("src/templates"))),
	)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	sheetId := r.URL.Query().Get("sheetId")
	if sheetId == "" {
		http.Error(w, "Missing sheetId parameter", http.StatusBadRequest)
		return
	}

	_, err := redisClient.Get(ctx, sheetId).Result()
	if err == redis.Nil {
		http.Error(w, "No locations found", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("src/templates/index.html"))
	tmpl.Execute(w, map[string]string{"SheetId": sheetId})
}

func loadLocationsHandler(w http.ResponseWriter, r *http.Request) {
	sheetId := r.URL.Query().Get("sheetId")
	if sheetId == "" {
		http.Error(w, "Missing sheetId parameter", http.StatusBadRequest)
		return
	}

	spreadsheetUrl := fmt.Sprintf(baseSpreadsheetUrl, sheetId)

	parsedLocations, err := locationHelpers.LoadLocations(spreadsheetUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	locationsJson, _ := json.Marshal(parsedLocations)
	redisClient.Set(ctx, sheetId, locationsJson, 0)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Locations located and cached"))
}

func locationsHandler(w http.ResponseWriter, r *http.Request) {
	sheetId := r.URL.Query().Get("sheetId")
	if sheetId == "" {
		http.Error(w, "Missing sheetId parameter", http.StatusBadRequest)
		return
	}

	locationsJson, err := redisClient.Get(r.Context(), sheetId).Bytes()
	if err == redis.Nil {
		http.Error(w, "No locations found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(locationsJson))
}

// func reloadLocationsHandler(w http.ResponseWriter, r *http.Request, setLocations func([]locationHelpers.Location), locationStore *LocationStore) {
// 	reloadedLocations, err := locationHelpers.LoadLocations()
// 	if err != nil {
// 		log.Printf("Failed to reload locations: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}

// 	currentLocations := locationStore.GetLocations()
// 	diff := locationHelpers.DiffLocations(currentLocations, reloadedLocations)

// 	setLocations(reloadedLocations)

// 	log.Printf("current locations: %d", len(currentLocations))
// 	log.Printf("Reloaded %d locations", len(reloadedLocations))
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(diff)
// }
