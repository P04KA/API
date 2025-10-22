package usecase

import (
	"context"

	"github.com/P04KA/API/internal/models"
	"github.com/P04KA/API/pkg/stats"
)

type UserProvider interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id string) error
}

type StatsUseCase interface {
	GetStats(ctx context.Context, period string) (*stats.UserStatsResponse, error)
}
