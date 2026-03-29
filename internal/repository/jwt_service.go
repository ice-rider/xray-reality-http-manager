package repository

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"xray_server/internal/domain"
)

type JWTService struct {
	secretKey    []byte
	tokenTimeout time.Duration
}

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey string, tokenTimeout time.Duration) *JWTService {
	return &JWTService{
		secretKey:    []byte(secretKey),
		tokenTimeout: tokenTimeout,
	}
}

func (s *JWTService) GenerateToken(user *domain.User) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(s.tokenTimeout)

	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, int64(s.tokenTimeout.Seconds()), nil
}

func (s *JWTService) ValidateToken(tokenString string) (*domain.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return &domain.User{
		ID:       claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
	}, nil
}
