package main

import (
	"encoding/json"
	"gomap/locationHelpers"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func initRouter(locationStore *LocationStore, setLocations func([]locationHelpers.Location)) {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)

	r.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
		locationsHandler(w, r, locationStore)
	})

	r.HandleFunc("/reload-locations", func(w http.ResponseWriter, r *http.Request) {
		reloadLocationsHandler(w, r, setLocations)
	}).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func locationsHandler(w http.ResponseWriter, r *http.Request, locationStore *LocationStore) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locationStore.GetLocations())
}

func reloadLocationsHandler(w http.ResponseWriter, r *http.Request, setLocations func([]locationHelpers.Location)) {
	reloadedLocations, err := locationHelpers.LoadLocations()
	if err != nil {
		log.Printf("Failed to reload locations: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	setLocations(reloadedLocations)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reloadedLocations)
}
