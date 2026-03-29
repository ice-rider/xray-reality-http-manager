package http

import (
	"encoding/json"
	"net/http"
	"xray_server/internal/usecase/auth"
)

type AuthHandler struct {
	loginUseCase *auth.LoginUseCase
}

func NewAuthHandler(loginUseCase *auth.LoginUseCase) *AuthHandler {
	return &AuthHandler{
		loginUseCase: loginUseCase,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	output, err := h.loginUseCase.Execute(auth.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		http.Error(w, `{"error": "invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      output.Token,
		"expires_in": output.ExpiresIn,
	})
}
