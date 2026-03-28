package env

import (
	"fmt"
	"os"
	"strings"
)

// Config содержит все необходимые переменные окружения
type Config struct {
	Mldsa65Sign   string
	Mldsa65Public string
	PrivateKey    string
	PublicKey     string
	ShortIdsRaw   string
	ServerIP      string
}

// Load загружает и валидирует переменные окружения
func Load() (*Config, error) {
	cfg := &Config{
		Mldsa65Sign:   os.Getenv("mldsa65_sign"),
		Mldsa65Public: os.Getenv("mldsa65_public"),
		PrivateKey:    os.Getenv("private_key"),
		PublicKey:     os.Getenv("public_key"),
		ShortIdsRaw:   os.Getenv("shorts_id"),
		ServerIP:      os.Getenv("server_ip"),
	}

	// Валидация
	var errs []string

	if cfg.Mldsa65Sign == "" {
		errs = append(errs, "mldsa65_sign не установлена")
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
		return nil, fmt.Errorf("ошибка валидации переменных окружения:\n  - %s",
			strings.Join(errs, "\n  - "))
	}

	return cfg, nil
}
