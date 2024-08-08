package main

import (
	"encoding/json"
	"gomap/locationHelpers"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func initRouter(locations []locationHelpers.Location) {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)

	r.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
		locationsHandler(w, r, locations)
	})

	r.HandleFunc("/reload-locations", func(w http.ResponseWriter, r *http.Request) {
		locations = reloadLocationsHandler(w, r)
	}).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func locationsHandler(w http.ResponseWriter, r *http.Request, locations []locationHelpers.Location) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func reloadLocationsHandler(w http.ResponseWriter, r *http.Request) []locationHelpers.Location {
	reloadedLocations, err := locationHelpers.LoadLocations()
	if err != nil {
		log.Printf("Failed to reload locations: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	locationsHandler(w, r, reloadedLocations)
	return reloadedLocations
}
