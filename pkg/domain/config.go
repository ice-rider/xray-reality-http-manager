package domain

type Config struct {
	Log       Log        `json:"log"`
	Inbounds  []Inbound  `json:"inbounds"`
	Outbounds []Outbound `json:"outbounds"`
	Api       Api        `json:"api"`
	Stats     Stats      `json:"stats"`
	Policy    Policy     `json:"policy"`
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
}

type Settings struct {
	Clients    []Client `json:"clients"`
	Decryption string   `json:"decryption"`
}

type Client struct {
	ID   string `json:"id"`
	Flow string `json:"flow"`
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
	Mldsa65Sign   string   `json:"mldsa65Sign"`
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
