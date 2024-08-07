package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
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
	err := fetchLocationsFromGoogleSheets()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func fetchLocationsFromGoogleSheets() error {
	resp, err := http.Get("https://docs.google.com/spreadsheets/d/e/2PACX-1vQcLVeU_Wg9kZzdTQyHovryufU-EBy7nf--9uKNs7q-lPj00Drs9Y038q2IsrjZ_Ha1Kl5dmWYygeMq/pub?gid=0&single=true&output=csv")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	var locs []Location

	if _, err := reader.Read(); err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		latRowIndex := 5
		longRowIndex := 6

		lat, _ := strconv.ParseFloat(record[latRowIndex], 64)
		long, _ := strconv.ParseFloat(record[longRowIndex], 64)
		locs = append(locs, Location{
			Name:       record[0],
			Address:    record[1],
			City:       record[2],
			State:      record[3],
			Country:    record[4],
			Latitude:   lat,
			Longitude:  long,
			IsCo404Loc: containsCo404(record[0]),
		})
	}
	locations = locs
	return nil
}

func init() {
	err := fetchLocationsFromGoogleSheets()
	if err != nil {
		log.Fatalf("Failed to fetch initial locations: %v", err)
	}
}

func containsCo404(s string) bool {
	lowerStr := strings.ToLower(s)
	return strings.Contains(lowerStr, "co404")
}
