package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"xray_server/internal/domain"

	stats "github.com/v2fly/v2ray-core/v5/app/stats/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StatsRepositorygRPC struct {
	conn   *grpc.ClientConn
	client stats.StatsServiceClient
	mu     sync.Mutex
}

func NewStatsRepositorygRPC(endpoint string) (*StatsRepositorygRPC, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Xray API: %w", err)
	}

	return &StatsRepositorygRPC{
		conn:   conn,
		client: stats.NewStatsServiceClient(conn),
	}, nil
}

func (r *StatsRepositorygRPC) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

func (r *StatsRepositorygRPC) GetTrafficStats(email string) (*domain.UserTraffic, error) {
	allStats, err := r.GetAllTrafficStats()
	if err != nil {
		return nil, err
	}
	for _, stat := range allStats {
		if stat.Email == email {
			return &stat, nil
		}
	}
	return nil, nil
}

func (r *StatsRepositorygRPC) GetAllTrafficStats() ([]domain.UserTraffic, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ctx := context.Background()
	req := &stats.QueryStatsRequest{
		Pattern: "user>>>",
	}

	resp, err := r.client.QueryStats(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to query stats: %w", err)
	}

	userStats := make(map[string]*domain.UserTraffic)
	for _, stat := range resp.GetStat() {
		name := stat.GetName()
		value := stat.GetValue()
		email := extractEmailFromStatName(name)
		if email == "" {
			continue
		}
		if _, exists := userStats[email]; !exists {
			userStats[email] = &domain.UserTraffic{
				Email: email,
			}
		}
		if isUplink(name) {
			userStats[email].Uplink = value
		} else if isDownlink(name) {
			userStats[email].Downlink = value
		}
	}

	result := make([]domain.UserTraffic, 0, len(userStats))
	for _, s := range userStats {
		result = append(result, *s)
	}
	return result, nil
}

func extractEmailFromStatName(name string) string {
	start := 5
	if len(name) <= start {
		return ""
	}
	end := findSubstring(name, ">>>traffic>>>", start)
	if end == -1 {
		return ""
	}
	return name[start:end]
}

func isUplink(name string) bool {
	return len(name) > 7 && name[len(name)-7:] == "uplink"
}

func isDownlink(name string) bool {
	return len(name) > 8 && name[len(name)-8:] == "downlink"
}

func findSubstring(s, substr string, start int) int {
	idx := -1
	if start+len(substr) <= len(s) {
		for i := start; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				idx = i
				break
			}
		}
	}
	return idx
}
