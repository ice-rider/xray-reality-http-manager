package app

import (
	"encoding/json"
	"os"
	"sync"

	"xray_server/pkg/domain"
	"xray_server/pkg/xray"
)

type ConfigService struct {
	config     *domain.Config
	configPath string
	serverIP   string
	mu         sync.RWMutex
}

type ConfigServiceOptions struct {
	Mldsa65Sign   string
	Mldsa65Public string
	PrivateKey    string
	PublicKey     string
	ShortIdsRaw   string
	ConfigPath    string
	ServerIP      string
}

func NewConfigService(opts ConfigServiceOptions) *ConfigService {
	shortIds := xray.ParseShortIds(opts.ShortIdsRaw)

	config := &domain.Config{
		Log: domain.Log{
			LogLevel: "warning",
		},
		Inbounds: []domain.Inbound{
			{
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
						ShortIds:      shortIds,
						Fingerprint:   "firefox",
						Mldsa65Sign:   opts.Mldsa65Sign,
						Mldsa65Public: opts.Mldsa65Public,
					},
				},
			},
		},
		Outbounds: []domain.Outbound{
			{
				Protocol: "freedom",
				Tag:      "direct",
			},
		},
		Api: domain.Api{
			Tag: "api",
			Services: []string{"HandlerService", "LoggerService", "StatsService"},
		},
		Stats: domain.Stats{},
		Policy: domain.Policy{
			Levels: map[string]domain.PolicyLevel0{
				"0": {
					StatsUserUplink: 	true,
					StatsUserDownlink: 	true,
				},
			},
		},
	}

	return &ConfigService{
		config:     config,
		configPath: opts.ConfigPath,	
		serverIP:   opts.ServerIP,
	}
}

func (s *ConfigService) AddClient(id, flow string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	client := domain.Client{
		ID:   id,
		Flow: flow,
	}

	s.config.Inbounds[0].Settings.Clients = append(s.config.Inbounds[0].Settings.Clients, client)

	return s.saveConfig()
}

func (s *ConfigService) GetConfig() *domain.Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

func (s *ConfigService) saveConfig() error {
	data, err := json.MarshalIndent(s.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.configPath, data, 0o644)
}

func (s *ConfigService) SaveConfig() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveConfig()
}

func (s *ConfigService) GetServerIP() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.serverIP
}

func (s *ConfigService) GetPublicKey() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config.Inbounds[0].StreamSettings.RealitySettings.PublicKey
}

func (s *ConfigService) GetMldsa65Public() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config.Inbounds[0].StreamSettings.RealitySettings.Mldsa65Public
}

func (s *ConfigService) GetShortIds() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config.Inbounds[0].StreamSettings.RealitySettings.ShortIds
}
