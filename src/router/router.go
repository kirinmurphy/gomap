package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

var templateDir = "src/templates"

func InitRouter(routerConfig RouterConfig) *mux.Router {
	r := mux.NewRouter()

	InitializeHomePageTemplates(templateDir)
	InitializeUpdateMapUITemplates(templateDir)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeRouteHandler(w, r, routerConfig)
	})

	r.HandleFunc("/updateMapUI", func(w http.ResponseWriter, r *http.Request) {
		updateMapUIHandler(w, r, routerConfig)
	}).Methods("POST")

	r.HandleFunc("/loadLocations", func(w http.ResponseWriter, r *http.Request) {
		loadLocationsRouteHandler(w, r, routerConfig)
	}).Methods("GET")

	r.HandleFunc("/getLocations", func(w http.ResponseWriter, r *http.Request) {
		getLocationsRouteHandler(w, r, routerConfig.RedisClient)
	})

	r.PathPrefix("/src/templates/").Handler(
		http.StripPrefix("/src/templates/", http.FileServer(http.Dir("src/templates"))),
	)

	return r
}
