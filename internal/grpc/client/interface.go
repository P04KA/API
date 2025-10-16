package client

import (
	"context"

	"github.com/P04KA/API/pkg/stats"
)

type StatsClient interface {
	GetUserStats(ctx context.Context, period string) (*stats.UserStatsResponse, error)
	Close() error
}
