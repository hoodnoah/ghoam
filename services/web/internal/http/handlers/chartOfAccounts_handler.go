package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/hoodnoah/ghoam/internal/services"
)

type ChartOfAccountsHandler struct {
	ChartOfAccountsService  *services.ChartOfAccountsService
	ChartOfAccountsTemplate *template.Template
}

func (h *ChartOfAccountsHandler) GetChart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "text/html")

	// fetch chart of accounts
	chart, err := h.ChartOfAccountsService.GetChartOfAccounts(ctx)
	if err != nil {
		log.Printf("failed to get chart of accounts with error %v", err)
		http.Error(w, "failed to get chart of accounts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// render full page
	if err := h.ChartOfAccountsTemplate.ExecuteTemplate(w, "chart", chart); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "render error: "+err.Error(), http.StatusInternalServerError)
	}

	// if it's an hx-swap, render only the fragment
	if r.Header.Get("HX-Request") == "true" {
		if err := h.ChartOfAccountsTemplate.ExecuteTemplate(w, "chartFragment", chart); err != nil {
			log.Printf("failed to render chart of accounts fragment with error :%v", err)
			http.Error(w, "Failed to render chart of accounts fragment:"+err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
