package http

import (
	"encoding/json"
	"net/http"

	"xray_server/internal/app"
	"xray_server/pkg/xray"
)

const (
	serverPort       = 443
	defaultFingerprint = "firefox"
	defaultSni       = "www.apple.com"
)

// ClientHandler обрабатывает HTTP запросы для управления клиентами
type ClientHandler struct {
	configService *app.ConfigService
}

// NewClientHandler создаёт новый хендлер
func NewClientHandler(configService *app.ConfigService) *ClientHandler {
	return &ClientHandler{
		configService: configService,
	}
}

// CreateClientRequest представляет запрос на создание клиента
type CreateClientRequest struct {
	Flow       string `json:"flow"`
	ClientName string `json:"client_name"`
}

// CreateClientResponse представляет ответ на создание клиента
type CreateClientResponse struct {
	ID      string `json:"id"`
	Flow    string `json:"flow"`
	Link    string `json:"link"`
}

// CreateClient обрабатывает POST /api/v1/client
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

	// Генерируем UUID для клиента
	id := xray.GenerateUUID()

	if err := h.configService.AddClient(id, req.Flow); err != nil {
		http.Error(w, "Failed to add client", http.StatusInternalServerError)
		return
	}

	// Получаем данные для генерации ссылки
	serverIP := h.configService.GetServerIP()
	publicKey := h.configService.GetPublicKey()
	mldsa65Public := h.configService.GetMldsa65Public()
	shortIds := h.configService.GetShortIds()

	// Используем первый ShortId или пустую строку
	shortId := ""
	if len(shortIds) > 0 {
		shortId = shortIds[0]
	}

	// Генерируем VLESS ссылку
	link := xray.GenerateVlessLink(xray.VlessLinkParams{
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
		ID:      id,
		Flow:    req.Flow,
		Link:    link,
	}

	json.NewEncoder(w).Encode(resp)
}
