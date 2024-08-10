package router

import (
	"context"
	"net/http"
	"text/template"

	"github.com/redis/go-redis/v9"
)

func homeRouteHandler(w http.ResponseWriter, r *http.Request, redisClient *redis.Client, ctx context.Context) {
	sheetId := r.URL.Query().Get("sheetId")
	if sheetId == "" {
		tmpl := template.Must(template.ParseFiles("src/templates/home.html"))
		tmpl.Execute(w, nil)
		return
	}

	_, err := redisClient.Get(ctx, sheetId).Result()
	if err == redis.Nil {
		http.Error(w, "Could not find spreadsheet id ${sheetId}", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("src/templates/map.html"))
	tmpl.Execute(w, map[string]string{"SheetId": sheetId})
}
