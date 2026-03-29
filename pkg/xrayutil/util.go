package xrayutil

import (
	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func ParseShortIds(raw string) []string {
	if raw == "" {
		return []string{}
	}

	result := []string{}
	start := 0
	for i, r := range raw {
		if r == ',' || r == ' ' {
			if start < i {
				result = append(result, raw[start:i])
			}
			start = i + 1
		}
	}
	if start < len(raw) {
		result = append(result, raw[start:])
	}

	return result
}
