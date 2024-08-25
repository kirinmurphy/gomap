package router

import (
	"context"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/redis/go-redis/v9"
)

var (
	homepageTemplate *template.Template
	mapPageTemplate  *template.Template
)

func InitializeHomePageTemplates(templateDir string) {
	homepageTemplate = template.Must(template.ParseFiles(filepath.Join(templateDir, "home.html")))
	mapPageTemplate = template.Must(template.ParseFiles(filepath.Join(templateDir, "map.html")))
}

func homeRouteHandler(w http.ResponseWriter, r *http.Request, redisClient RedisClientInterface, ctx context.Context) {
	sheetId := r.URL.Query().Get("sheetId")
	if sheetId == "" {
		demoFlag := r.URL.Query().Get("demo")
		data := map[string]interface{}{
			"Demo": demoFlag == "true",
		}

		homepageTemplate.Execute(w, data)
		return
	}

	_, err := redisClient.Get(ctx, sheetId).Result()
	if err == redis.Nil {
		http.Error(w, "Could not find spreadsheet id ${sheetId}", http.StatusNotFound)
		return
	}

	mapPageTemplate.Execute(w, map[string]string{"SheetId": sheetId})
}
