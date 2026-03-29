package http

import (
	"encoding/json"
	"net/http"
	"xray_server/internal/usecase/config"
)

type ServerConfigHandler struct {
	configUseCase *config.ConfigUseCase
}

func NewServerConfigHandler(configUseCase *config.ConfigUseCase) *ServerConfigHandler {
	return &ServerConfigHandler{
		configUseCase: configUseCase,
	}
}

type ServerConfigResponse struct {
	ServerIP      string   `json:"server_ip"`
	PublicKey     string   `json:"public_key"`
	Mldsa65Public string   `json:"mldsa65_public"`
	ShortIds      []string `json:"short_ids"`
	Port          int      `json:"port"`
	Sni           string   `json:"sni"`
	Fingerprint   string   `json:"fingerprint"`
}

func (h *ServerConfigHandler) GetServerConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := ServerConfigResponse{
		ServerIP:      h.configUseCase.GetServerIP(),
		PublicKey:     h.configUseCase.GetPublicKey(),
		Mldsa65Public: h.configUseCase.GetMldsa65Public(),
		ShortIds:      h.configUseCase.GetShortIds(),
		Port:          443,
		Sni:           "www.apple.com",
		Fingerprint:   "firefox",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(config)
}
