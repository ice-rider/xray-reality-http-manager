package config

import (
	"encoding/json"
	"os"
	"sync"
	"xray_server/internal/domain"
)

type ConfigUseCase struct {
	config     *domain.Config
	configPath string
	serverIP   string
	mu         sync.RWMutex
}

type ConfigOptions struct {
	Mldsa65Seed   string
	Mldsa65Public string
	PrivateKey    string
	PublicKey     string
	ShortIdsRaw   string
	ConfigPath    string
	ServerIP      string
}

func NewConfigUseCase(opts ConfigOptions) *ConfigUseCase {
	config := &domain.Config{
		Log: domain.Log{LogLevel: "warning"},
		Inbounds: []domain.InboundBase{
			domain.Inbound{
				Listen:   "0.0.0.0",
				Port:     443,
				Protocol: "vless",
				Settings: domain.Settings{
					Clients:    []domain.Client{},
					Decryption: "none",
				},
				StreamSettings: domain.StreamSettings{
					Network:  "tcp",
					Security: "reality",
					RealitySettings: domain.RealitySettings{
						Show:          false,
						Dest:          "www.apple.com:443",
						Xver:          0,
						ServerNames:   []string{"www.apple.com"},
						PrivateKey:    opts.PrivateKey,
						PublicKey:     opts.PublicKey,
						ShortIds:      parseShortIds(opts.ShortIdsRaw),
						Fingerprint:   "firefox",
						Mldsa65Seed:   opts.Mldsa65Seed,
						Mldsa65Public: opts.Mldsa65Public,
					},
				},
			},
			domain.ApiInbound{
				Listen:   "0.0.0.0",
				Port:     54321,
				Protocol: "dokodemo-door",
				Settings: domain.ApiInboundSettings{Address: "127.0.0.1"},
				Tag:      "api",
			},
		},
		Outbounds: []domain.Outbound{{Protocol: "freedom", Tag: "direct"}},
		Api:       domain.Api{Tag: "api", Services: []string{"HandlerService", "LoggerService", "StatsService"}},
		Stats:     domain.Stats{},
		Policy: domain.Policy{
			Levels: map[string]domain.PolicyLevel0{
				"0": {StatsUserUplink: true, StatsUserDownlink: true},
			},
		},
		Routing: domain.Routing{
			Rules: []domain.RoutingRule{
				{
					Type:        "field", // без этого xray молча игнорирует правило
					InboundTag:  []string{"api"},
					OutboundTag: "api",
				},
			},
		},
	}

	return &ConfigUseCase{
		config:     config,
		configPath: opts.ConfigPath,
		serverIP:   opts.ServerIP,
	}
}

type AddClientInput struct {
	ID    string
	Flow  string
	Email string
}

func (uc *ConfigUseCase) AddClient(input AddClientInput) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	client := domain.Client{ID: input.ID, Flow: input.Flow, Email: input.Email}
	if inbound, ok := uc.config.Inbounds[0].(domain.Inbound); ok {
		inbound.Settings.Clients = append(inbound.Settings.Clients, client)
		uc.config.Inbounds[0] = inbound
	}
	return uc.saveConfig()
}

func (uc *ConfigUseCase) SaveConfig() error {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	return uc.saveConfig()
}

func (uc *ConfigUseCase) saveConfig() error {
	data, err := json.MarshalIndent(uc.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(uc.configPath, data, 0o644)
}

func (uc *ConfigUseCase) GetServerIP() string {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	return uc.serverIP
}

func (uc *ConfigUseCase) GetPublicKey() string {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if inbound, ok := uc.config.Inbounds[0].(domain.Inbound); ok {
		return inbound.StreamSettings.RealitySettings.PublicKey
	}
	return ""
}

func (uc *ConfigUseCase) GetMldsa65Public() string {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if inbound, ok := uc.config.Inbounds[0].(domain.Inbound); ok {
		return inbound.StreamSettings.RealitySettings.Mldsa65Public
	}
	return ""
}

func (uc *ConfigUseCase) GetShortIds() []string {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if inbound, ok := uc.config.Inbounds[0].(domain.Inbound); ok {
		return inbound.StreamSettings.RealitySettings.ShortIds
	}
	return []string{}
}

func (uc *ConfigUseCase) GetClients() []domain.Client {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if inbound, ok := uc.config.Inbounds[0].(domain.Inbound); ok {
		return inbound.Settings.Clients
	}
	return []domain.Client{}
}

func parseShortIds(raw string) []string {
	if raw == "" {
		return []string{}
	}
	var result []string
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
