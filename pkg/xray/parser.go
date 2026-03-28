package xray

import "strings"

// ParseShortIds парсит строку с Short IDs в слайс
func ParseShortIds(raw string) []string {
	if raw == "" {
		return []string{}
	}

	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' '
	})

	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
