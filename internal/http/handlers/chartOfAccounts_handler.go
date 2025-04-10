package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hoodnoah/ghoam/internal/services"
)

type ChartOfAccountsHandler struct {
	ChartOfAccountsService *services.ChartOfAccountsService
}

func (h *ChartOfAccountsHandler) GetChart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	chart, err := h.ChartOfAccountsService.GetChartOfAccounts(ctx)
	if err != nil {
		http.Error(w, "failed to get chart of accounts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chart); err != nil {
		http.Error(w, "failed to encode chart to json: "+err.Error(), http.StatusInternalServerError)
	}
}
