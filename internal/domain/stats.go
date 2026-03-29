package domain

type UserTraffic struct {
	Email    string `json:"email"`
	Uplink   int64  `json:"uplink"`
	Downlink int64  `json:"downlink"`
}

type StatsRepository interface {
	GetTrafficStats(email string) (*UserTraffic, error)
	GetAllTrafficStats() ([]UserTraffic, error)
}
