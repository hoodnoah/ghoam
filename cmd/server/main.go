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
	tmpl, err := template.ParseGlob(filepath.Join("templates", "*.gohtml"))
	if err != nil {
		log.Fatalf("failed to parse templates with error %v", err)
	}

	// Create the handler for the Chart of Accounts endpoint
	chartHandler := &handlers.ChartOfAccountsHandler{
		ChartOfAccountsService:  &chartService,
		ChartOfAccountsTemplate: tmpl,
	}

	// Set up routes: the index page and the chart endpoint for HTMX
	// index handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{"Title": "Home "}
		if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
			http.Error(w, "rendering index failed:"+err.Error(), http.StatusInternalServerError)
		}
	})

	// chart of accounts handler
	http.HandleFunc("/chart", chartHandler.GetChart)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
