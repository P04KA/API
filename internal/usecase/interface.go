package usecase

import (
	"context"

	"github.com/P04KA/API/internal/models"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, user models.User) (string, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id string) error
}
