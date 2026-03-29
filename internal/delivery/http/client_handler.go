package http

import (
	"encoding/json"
	"net/http"
	"xray_server/internal/usecase/config"
	"xray_server/pkg/xrayutil"
)

const (
	serverPort       = 443
	defaultFingerprint = "firefox"
	defaultSni       = "www.apple.com"
)

type ClientHandler struct {
	configUseCase *config.ConfigUseCase
}

func NewClientHandler(configUseCase *config.ConfigUseCase) *ClientHandler {
	return &ClientHandler{
		configUseCase: configUseCase,
	}
}

type CreateClientRequest struct {
	Flow       string `json:"flow"`
	ClientName string `json:"client_name"`
	Email      string `json:"email,omitempty"`
}

type CreateClientResponse struct {
	ID   string `json:"id"`
	Flow string `json:"flow"`
	Link string `json:"link"`
}

type GetClientsResponse struct {
	Clients []ClientInfo `json:"clients"`
}

type ClientInfo struct {
	ID    string `json:"id"`
	Flow  string `json:"flow"`
	Email string `json:"email"`
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := xrayutil.GenerateUUID()

	email := req.Email
	if email == "" && req.ClientName != "" {
		email = req.ClientName + "@local"
	}

	if err := h.configUseCase.AddClient(config.AddClientInput{
		ID:    id,
		Flow:  req.Flow,
		Email: email,
	}); err != nil {
		http.Error(w, "Failed to add client", http.StatusInternalServerError)
		return
	}

	serverIP := h.configUseCase.GetServerIP()
	publicKey := h.configUseCase.GetPublicKey()
	mldsa65Public := h.configUseCase.GetMldsa65Public()
	shortIds := h.configUseCase.GetShortIds()

	shortId := ""
	if len(shortIds) > 0 {
		shortId = shortIds[0]
	}

	link := xrayutil.GenerateVlessLink(xrayutil.VlessLinkParams{
		UUID:        id,
		ServerIP:    serverIP,
		Port:        serverPort,
		Flow:        req.Flow,
		PublicKey:   publicKey,
		ShortId:     shortId,
		Mldsa65Pqv:  mldsa65Public,
		Fingerprint: defaultFingerprint,
		Sni:         defaultSni,
		ClientName:  req.ClientName,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := CreateClientResponse{
		ID:   id,
		Flow: req.Flow,
		Link: link,
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *ClientHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clients := h.configUseCase.GetClients()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := GetClientsResponse{
		Clients: make([]ClientInfo, len(clients)),
	}

	for i, client := range clients {
		resp.Clients[i] = ClientInfo{
			ID:    client.ID,
			Flow:  client.Flow,
			Email: client.Email,
		}
	}

	json.NewEncoder(w).Encode(resp)
}
