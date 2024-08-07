package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"

	"gopro/locationHelpers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var locations = []locationHelpers.Location{}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/locations", locationsHandler)
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

func init() {
	godotenv.Load()
	spreadsheetUrl := os.Getenv("GOOGLE_SHEET_CSV_URL")

	csvStream, errChan := locationHelpers.FetchLocationData(spreadsheetUrl)

	parsedLocations, err := locationHelpers.ParseLocations(csvStream)
	if err != nil {
		log.Fatalf("Failed to parse locations: %v", err)
	}

	if err := <-errChan; err != nil {
		log.Fatalf("Failed to fetch CSV data: %v", err)
	}

	locations = parsedLocations
}
