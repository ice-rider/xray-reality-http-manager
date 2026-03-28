package xray

import (
	"fmt"
	"net/url"
)

// VlessLinkParams параметры для генерации VLESS ссылки
type VlessLinkParams struct {
	UUID        string
	ServerIP    string
	Port        int
	Flow        string
	PublicKey   string
	ShortId     string
	Mldsa65Pqv  string
	Fingerprint string
	Sni         string
	ClientName  string
}

// GenerateVlessLink генерирует VLESS ссылку в формате:
// vless://uuid@ip:port?encryption=none&flow=...&security=reality&sni=...&fp=...&pbk=...&sid=...&pqv=...#name
func GenerateVlessLink(params VlessLinkParams) string {
	query := url.Values{}
	query.Set("encryption", "none")
	query.Set("flow", params.Flow)
	query.Set("security", "reality")
	query.Set("sni", params.Sni)
	query.Set("fp", params.Fingerprint)
	query.Set("pbk", params.PublicKey)
	query.Set("sid", params.ShortId)
	query.Set("pqv", params.Mldsa65Pqv)

	u := &url.URL{
		Scheme:   "vless",
		User:     url.User(params.UUID),
		Host:     fmt.Sprintf("%s:%d", params.ServerIP, params.Port),
		RawQuery: query.Encode(),
	}

	if params.ClientName != "" {
		u.Fragment = params.ClientName
	}

	return u.String()
}
