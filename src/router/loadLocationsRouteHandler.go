package router

import (
	"encoding/json"
	"net/http"
)

func loadLocationsRouteHandler(w http.ResponseWriter, r *http.Request, routerConfig RouterConfig) {
	sheetId := r.URL.Query().Get("sheetId")
	if sheetId == "" {
		http.Error(w, "Missing sheetId parameter", http.StatusBadRequest)
		return
	}

	err := processLocations(sheetId, routerConfig)
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
