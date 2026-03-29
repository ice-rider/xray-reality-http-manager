package http

import (
	"encoding/json"
	"net/http"
	"xray_server/internal/usecase/config"
)

type Mldsa65Handler struct {
	configUseCase *config.ConfigUseCase
}

func NewMldsa65Handler(configUseCase *config.ConfigUseCase) *Mldsa65Handler {
	return &Mldsa65Handler{
		configUseCase: configUseCase,
	}
}

type Mldsa65PublicResponse struct {
	Mldsa65Public string `json:"mldsa65_public"`
}

func (h *Mldsa65Handler) GetMldsa65Public(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mldsa65Public := h.configUseCase.GetMldsa65Public()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := Mldsa65PublicResponse{
		Mldsa65Public: mldsa65Public,
	}

	json.NewEncoder(w).Encode(resp)
}
