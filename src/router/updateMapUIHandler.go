package router

import (
	"context"
	"html/template"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var (
	errorTemplate   = template.Must(template.ParseFiles("src/templates/loadLocationsError.html"))
	successTemplate = template.Must(template.ParseFiles("src/templates/loadlLocationsSuccess.html"))
)

func updateMapUIHandler(w http.ResponseWriter, r *http.Request, redisClient *redis.Client, ctx context.Context) {
	w.Header().Set("Content-Type", "text/html")

	sheetId := r.FormValue("sheetId")
	if sheetId == "" {
		renderErrorHTML(w, "Missing sheetId parameter", http.StatusBadRequest)
		return
	}

	err := processLocations(sheetId, redisClient, ctx)
	if err != nil {
		renderErrorHTML(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = successTemplate.Execute(w, map[string]string{"SheetId": sheetId})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func renderErrorHTML(w http.ResponseWriter, errMsg string, statusCode int) {
	w.WriteHeader(statusCode)
	err := errorTemplate.Execute(w, map[string]string{"Error": errMsg})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
