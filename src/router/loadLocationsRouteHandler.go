package router

import (
	"context"
	"encoding/json"
	"fmt"
	"gomap/src/locationManager"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

const baseSpreadsheetUrl = "https://docs.google.com/spreadsheets/d/e/%s/pub?gid=0&single=true&output=csv"

func loadLocationsRouteHandler(w http.ResponseWriter, r *http.Request, redisClient *redis.Client, ctx context.Context) {
	var sheetId string

	switch r.Method {
	case http.MethodGet:
		sheetId = r.URL.Query().Get("sheetId")
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		sheetId = r.FormValue("sheetId")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("SheetId: %s", sheetId)
	if sheetId == "" {
		http.Error(w, "Missing sheetId parameter", http.StatusBadRequest)
		return
	}

	spreadsheetUrl := fmt.Sprintf(baseSpreadsheetUrl, sheetId)

	parsedLocations, err := locationManager.LoadLocations(spreadsheetUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	locationsJson, err := json.Marshal(parsedLocations)
	if err != nil {
		http.Error(w, "Failed to marshal locations", http.StatusInternalServerError)
		return
	}

	err = redisClient.Set(ctx, sheetId, locationsJson, 0).Err()
	if err != nil {
		http.Error(w, "Failed to cache locations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"success": "true",
		"message": "Locations loaded and cached",
		"sheetId": sheetId,
	})
}
