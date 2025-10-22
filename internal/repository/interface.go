package repository

import (
	"context"
	"time"

	"github.com/P04KA/API/internal/models"
)

type UserProvider interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id string) error
	CountUserCreated(ctx context.Context, time time.Time) (int64, error)
	CountUserUpdated(ctx context.Context, time time.Time) (int64, error)
	CountUserDeleted(ctx context.Context, time time.Time) (int64, error)
}
