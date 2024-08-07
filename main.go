package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Location struct {
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	Country    string  `json:"country"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	IsCo404Loc bool    `json:"isCo404Loc"`
}

var locations = []Location{}

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
	url := os.Getenv("GOOGLE_SHEET_CSV_URL")

	err := fetchLocationsFromGoogleSheets(url)
	if err != nil {
		log.Fatalf("Failed to fetch initial locations: %v", err)
	}
}
