package handlers

import (
	"html/template"
	"net/http"

	"github.com/hoodnoah/ghoam/internal/services"
)

type ChartOfAccountsHandler struct {
	ChartOfAccountsService  *services.ChartOfAccountsService
	ChartOfAccountsTemplate *template.Template
}

func (h *ChartOfAccountsHandler) GetChart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	chart, err := h.ChartOfAccountsService.GetChartOfAccounts(ctx)
	if err != nil {
		http.Error(w, "failed to get chart of accounts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	// render the chart using the chart.tmpl template
	if err := h.ChartOfAccountsTemplate.Execute(w, chart); err != nil {
		http.Error(w, "Failed to render chart: "+err.Error(), http.StatusInternalServerError)
	}
}
