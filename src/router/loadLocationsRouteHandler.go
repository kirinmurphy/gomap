package router

import (
	"context"
	"encoding/json"
	"fmt"
	"gomap/src/locationHelpers"
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
