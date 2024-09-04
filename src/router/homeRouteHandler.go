package router

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/redis/go-redis/v9"
)

type HomePageTemplateData struct {
	Demo bool
}

type MapPageTemplateData struct {
	SheetId string
}

var (
	homepageTemplate *template.Template
	mapPageTemplate  *template.Template
)

func InitializeHomePageTemplates(templateDir string) {
	homepageTemplate = template.Must(template.ParseFiles(filepath.Join(templateDir, "home.html")))
	mapPageTemplate = template.Must(template.ParseFiles(filepath.Join(templateDir, "map.html")))
}

func homeRouteHandler(w http.ResponseWriter, r *http.Request, config RouterConfig) {
	query := r.URL.Query()
	sheetId := query.Get("sheetId")

	if sheetId == "" {
		sheetId = query.Get("sheetID")
	}

	if sheetId == "" {
		demoFlag := r.URL.Query().Get("demo")

		homepageTemplate.Execute(w, HomePageTemplateData{
			Demo: demoFlag == "true",
		})
		return
	}

	_, err := config.RedisClient.Get(config.Ctx, sheetId).Result()
	if err == redis.Nil {
		http.Error(w, "Could not find spreadsheet id ${sheetId}", http.StatusNotFound)
		return
	}

	mapPageTemplate.Execute(w, MapPageTemplateData{
		SheetId: sheetId,
	})
}
