package router

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	errorTemplate   *template.Template
	successTemplate *template.Template
)

func InitializeUpdateMapUITemplates(templateDir string) {
	errorTemplate = template.Must(template.ParseFiles(filepath.Join(templateDir, "loadLocationsError.html")))
	successTemplate = template.Must(template.ParseFiles(filepath.Join(templateDir, "loadlLocationsSuccess.html")))
}

func updateMapUIHandler(w http.ResponseWriter, r *http.Request, routerConfig RouterConfig) {
	w.Header().Set("Content-Type", "text/html")

	sheetId := r.FormValue("sheetId")
	if sheetId == "" {
		renderErrorHTML(w, "Missing sheetId parameter", http.StatusBadRequest)
		return
	}

	err := processLocations(sheetId, routerConfig)
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
