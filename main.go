package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"gopro/locationHelpers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var locations = []locationHelpers.Location{}

func init() {
	godotenv.Load()
	if err := loadLocations(); err != nil {
		log.Fatalf("failed to load initial locations: %v", err)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/locations", locationsHandler)
	r.HandleFunc("/reload-locations", reloadLocationsHandler).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func locationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func reloadLocationsHandler(w http.ResponseWriter, r *http.Request) {
	if err := loadLocations(); err != nil {
		log.Printf("Failed to reload locations: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadLocations() error {
	spreadsheetUrl := os.Getenv("GOOGLE_SHEET_CSV_URL")

	csvStream, errChan := locationHelpers.FetchLocationData(spreadsheetUrl)

	parsedLocations, err := locationHelpers.ParseLocations(csvStream)
	if err != nil {
		return fmt.Errorf("failed to parse locations: %w", err)
	}

	if err := <-errChan; err != nil {
		return fmt.Errorf("failed to fetch CSV data: %w", err)
	}

	locations = parsedLocations
	return nil
}
