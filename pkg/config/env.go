package config

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"
)

type EnvConfig struct {
	Mldsa65Seed     string
	Mldsa65Public   string
	PrivateKey      string
	PublicKey       string
	ShortIdsRaw     string
	ServerIP        string
	JWTSecret       string
	AdminUsername   string
	AdminPassword   string
}

func LoadEnv() (*EnvConfig, error) {
	cfg := &EnvConfig{
		Mldsa65Seed:     os.Getenv("mldsa65_seed"),
		Mldsa65Public:   os.Getenv("mldsa65_public"),
		PrivateKey:      os.Getenv("private_key"),
		PublicKey:       os.Getenv("public_key"),
		ShortIdsRaw:     os.Getenv("shorts_id"),
		ServerIP:        os.Getenv("server_ip"),
		JWTSecret:       os.Getenv("jwt_secret"),
		AdminUsername:   os.Getenv("admin_username"),
		AdminPassword:   os.Getenv("admin_password"),
	}

	var errs []string

	if cfg.Mldsa65Seed == "" {
		errs = append(errs, "mldsa65_seed не установлена")
	}
	if cfg.Mldsa65Public == "" {
		errs = append(errs, "mldsa65_public не установлена")
	}
	if cfg.PrivateKey == "" {
		errs = append(errs, "private_key не установлена")
	}
	if cfg.PublicKey == "" {
		errs = append(errs, "public_key не установлена")
	}
	if cfg.ShortIdsRaw == "" {
		errs = append(errs, "shorts_id не установлена")
	}
	if cfg.ServerIP == "" {
		errs = append(errs, "server_ip не установлена")
	}

	if len(errs) > 0 {
		return nil, &EnvError{errs: errs}
	}

	if cfg.AdminUsername == "" {
		cfg.AdminUsername = "admin"
	}
	if cfg.AdminPassword == "" {
		cfg.AdminPassword = "admin"
	}

	if cfg.JWTSecret == "" {
		secret, err := generateRandomSecret()
		if err != nil {
			return nil, err
		}
		cfg.JWTSecret = secret
	}

	return cfg, nil
}

type EnvError struct {
	errs []string
}

func (e *EnvError) Error() string {
	var sb strings.Builder
	sb.WriteString("ошибка валидации переменных окружения:\n  - ")
	for i, err := range e.errs {
		if i > 0 {
			sb.WriteString("\n  - ")
		}
		sb.WriteString(err)
	}
	return sb.String()
}

func generateRandomSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
