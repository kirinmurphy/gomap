package router

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
)

const baseSpreadsheetUrl = "https://docs.google.com/spreadsheets/d/e/%s/pub?gid=0&single=true&output=csv"

func loadLocationsRouteHandler(w http.ResponseWriter, r *http.Request, redisClient *redis.Client, ctx context.Context) {
	sheetId := r.URL.Query().Get("sheetId")
	if sheetId == "" {
		http.Error(w, "Missing sheetId parameter", http.StatusBadRequest)
		return
	}

	err := processLocations(sheetId, redisClient, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
