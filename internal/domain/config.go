package domain

type InboundBase interface {
	GetTag() string
}

type Config struct {
	Log       Log           `json:"log"`
	Inbounds  []InboundBase `json:"inbounds"`
	Outbounds []Outbound    `json:"outbounds"`
	Api       Api           `json:"api"`
	Stats     Stats         `json:"stats"`
	Policy    Policy        `json:"policy"`
	Routing   Routing       `json:"routing"`
}

type Routing struct {
	Rules []RoutingRule `json:"rules"`
}

type RoutingRule struct {
	Type        string   `json:"type"`
	InboundTag  []string `json:"inboundTag"`
	OutboundTag string   `json:"outboundTag"`
}

type ApiInbound struct {
	Listen   string             `json:"listen"`
	Port     int                `json:"port"`
	Protocol string             `json:"protocol"`
	Settings ApiInboundSettings `json:"settings"`
	Tag      string             `json:"tag"`
}

func (a ApiInbound) GetTag() string {
	return a.Tag
}

type ApiInboundSettings struct {
	Address string `json:"address"`
}

type Log struct {
	LogLevel string `json:"loglevel"`
}

type Inbound struct {
	Listen         string         `json:"listen"`
	Port           int            `json:"port"`
	Protocol       string         `json:"protocol"`
	Settings       Settings       `json:"settings"`
	StreamSettings StreamSettings `json:"streamSettings"`
	Tag            string         `json:"tag,omitempty"`
}

func (i Inbound) GetTag() string {
	return i.Tag
}

type Settings struct {
	Clients    []Client `json:"clients"`
	Decryption string   `json:"decryption"`
}

type Client struct {
	ID    string `json:"id"`
	Flow  string `json:"flow"`
	Email string `json:"email,omitempty"`
}

type StreamSettings struct {
	Network         string          `json:"network"`
	Security        string          `json:"security"`
	RealitySettings RealitySettings `json:"realitySettings"`
}

type RealitySettings struct {
	Show          bool     `json:"show"`
	Dest          string   `json:"dest"`
	Xver          int      `json:"xver"`
	ServerNames   []string `json:"serverNames"`
	PrivateKey    string   `json:"privateKey"`
	PublicKey     string   `json:"publicKey"`
	ShortIds      []string `json:"shortIds"`
	Fingerprint   string   `json:"fingerprint"`
	Mldsa65Seed   string   `json:"mldsa65Seed"`
	Mldsa65Public string   `json:"mldsa65Public"`
}

type Outbound struct {
	Protocol string `json:"protocol"`
	Tag      string `json:"tag"`
}

type Api struct {
	Tag      string   `json:"tag"`
	Services []string `json:"services"`
}

type Stats struct{}

type Policy struct {
	Levels map[string]PolicyLevel0 `json:"levels"`
}

type PolicyLevel0 struct {
	StatsUserUplink   bool `json:"statsUserUplink"`
	StatsUserDownlink bool `json:"statsUserDownlink"`
}

type ConfigRepository interface {
	Save(config *Config) error
	AddClient(client Client) error
	GetClients() []Client
	GetServerIP() string
	GetPublicKey() string
	GetMldsa65Public() string
	GetShortIds() []string
}
