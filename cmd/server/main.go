package server

import (
	// std
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	// internal
	"github.com/hoodnoah/ghoam/internal/http/handlers"
	"github.com/hoodnoah/ghoam/internal/persistence/sqlite"
	"github.com/hoodnoah/ghoam/internal/services"
)

func Execute(repos *sqlite.Repositories) {
	// Instantiate the ChartOfAccountsService using the SQLite repositories
	chartService := services.ChartOfAccountsService{
		AccountRepo:      repos.Accounts,
		AccountGroupRepo: repos.AccountGroups,
	}

	// Parse templates from the templates/ folder
	tmpl, err := template.ParseGlob(filepath.Join("templates", "*.tmpl"))
	if err != nil {
		log.Fatalf("failed to parse templates with error %v", err)
	}

	// Get the chart, index templates
	chartTmpl := tmpl.Lookup("chart.tmpl")
	if chartTmpl == nil {
		log.Fatalf("chart.tmpl not found")
	}

	indexTmpl := tmpl.Lookup("index.tmpl")
	if indexTmpl == nil {
		log.Fatalf("index.tmpl not found")
	}

	// Create the handler for the Chart of Accounts endpoint
	chartHandler := &handlers.ChartOfAccountsHandler{
		ChartOfAccountsService:  &chartService,
		ChartOfAccountsTemplate: chartTmpl,
	}

	// Set up routes: the index page and the chart endpoint for HTMX
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := indexTmpl.Execute(w, nil); err != nil {
			http.Error(w, "failed to render index: "+err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("/chart", chartHandler.GetChart)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
