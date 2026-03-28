package xray

import "github.com/google/uuid"

// GenerateUUID генерирует новый UUID для клиента
func GenerateUUID() string {
	return uuid.New().String()
}
