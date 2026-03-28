package env

import (
	"fmt"
	"os"
	"strings"
)

type EnvConfig struct {
	Mldsa65Sign   string
	Mldsa65Public string
	PrivateKey    string
	PublicKey     string
	ShortIdsRaw string
	ServerIP string
}

type EnvValidationError struct {
	Field   string
	Message string
}

func (e *EnvValidationError) Error() string {
	return fmt.Sprintf("invalid env variable %s: %s", e.Field, e.Message)
}

func Load() (*EnvConfig, error) {
	config := loadEnv()

	if err := validate(config); err != nil {
		return nil, err
	}

	return config, nil
}

func loadEnv() *EnvConfig {
	config := &EnvConfig{
		Mldsa65Sign:   os.Getenv("mldsa65_sign"),
		Mldsa65Public: os.Getenv("mldsa65_public"),
		PrivateKey:    os.Getenv("private_key"),
		PublicKey:     os.Getenv("public_key"),
		ShortIdsRaw:   os.Getenv("shorts_id"),
		ServerIP:      os.Getenv("server_ip"),
	}

	if config.ShortIdsRaw != "" {
		config.ShortIdsRaw = strings.Join(strings.Fields(strings.ReplaceAll(config.ShortIdsRaw, ",", " ")), " ")
	}

	return config
}

func validate(config *EnvConfig) error {
	if config.Mldsa65Sign == "" {
		return &EnvValidationError{
			Field:   "mldsa65_sign",
			Message: "variable is required but not set",
		}
	}

	if config.Mldsa65Public == "" {
		return &EnvValidationError{
			Field:   "mldsa65_public",
			Message: "variable is required but not set",
		}
	}

	if config.PrivateKey == "" {
		return &EnvValidationError{
			Field:   "private_key",
			Message: "variable is required but not set",
		}
	}

	if config.PublicKey == "" {
		return &EnvValidationError{
			Field:   "public_key",
			Message: "variable is required but not set",
		}
	}

	return nil
}
