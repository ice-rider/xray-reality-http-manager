package http

import (
	"encoding/json"
	"net/http"
	"xray_server/internal/usecase/stats"
)

type StatsHandler struct {
	statsUseCase *stats.StatsUseCase
}

func NewStatsHandler(statsUseCase *stats.StatsUseCase) *StatsHandler {
	return &StatsHandler{
		statsUseCase: statsUseCase,
	}
}

type GetClientsStatsResponse struct {
	Clients []stats.ClientTrafficStats `json:"clients"`
}

func (h *StatsHandler) GetClientsStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.statsUseCase.GetAllClientsStats()
	if err != nil {
		http.Error(w, "Failed to get stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := GetClientsStatsResponse{
		Clients: stats,
	}

	json.NewEncoder(w).Encode(resp)
}

type GetClientStatsResponse struct {
	Client *stats.ClientTrafficStats `json:"client"`
}

func (h *StatsHandler) GetClientStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parameter is required", http.StatusBadRequest)
		return
	}

	stats, err := h.statsUseCase.GetClientStats(email)
	if err != nil {
		http.Error(w, "Failed to get stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if stats == nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := GetClientStatsResponse{
		Client: stats,
	}

	json.NewEncoder(w).Encode(resp)
}
