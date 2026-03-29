package stats

import (
	"sync"
	"xray_server/internal/domain"
)

type ClientTrafficStats struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Flow     string `json:"flow"`
	Uplink   int64  `json:"uplink"`
	Downlink int64  `json:"downlink"`
	Total    int64  `json:"total"`
}

type StatsUseCase struct {
	statsRepo domain.StatsRepository
	mu        sync.RWMutex
}

func NewStatsUseCase(statsRepo domain.StatsRepository) *StatsUseCase {
	return &StatsUseCase{
		statsRepo: statsRepo,
	}
}

func (uc *StatsUseCase) GetAllClientsStats() ([]ClientTrafficStats, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	trafficStats, err := uc.statsRepo.GetAllTrafficStats()
	if err != nil {
		return nil, err
	}

	result := make([]ClientTrafficStats, 0, len(trafficStats))
	for _, traffic := range trafficStats {
		stats := ClientTrafficStats{
			Email:    traffic.Email,
			Uplink:   traffic.Uplink,
			Downlink: traffic.Downlink,
			Total:    traffic.Uplink + traffic.Downlink,
		}
		result = append(result, stats)
	}

	return result, nil
}

func (uc *StatsUseCase) GetClientStats(email string) (*ClientTrafficStats, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	traffic, err := uc.statsRepo.GetTrafficStats(email)
	if err != nil {
		return nil, err
	}

	if traffic == nil {
		return nil, nil
	}

	stats := &ClientTrafficStats{
		Email:    traffic.Email,
		Uplink:   traffic.Uplink,
		Downlink: traffic.Downlink,
		Total:    traffic.Uplink + traffic.Downlink,
	}

	return stats, nil
}
