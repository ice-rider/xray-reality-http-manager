package http

import (
	"context"
	"net/http"
	"strings"
	"xray_server/internal/domain"
	"xray_server/internal/repository"
)

type contextKey string

const userContextKey contextKey = "user"

type AuthMiddleware struct {
	jwtService *repository.JWTService
}

func NewAuthMiddleware(jwtService *repository.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error": "invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		user, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, `{"error": "invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	return user, ok
}
