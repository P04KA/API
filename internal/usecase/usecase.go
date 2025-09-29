package usecase

import (
	"context"

	"github.com/P04KA/API/internal/models"
	"github.com/P04KA/API/internal/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UseCase struct {
	userRepo repository.UserProvider
}

func New(userRepo repository.UserProvider) *UseCase {
	return &UseCase{
		userRepo: userRepo,
	}
}

func (u *UseCase) CreateUser(ctx context.Context, user models.User) (string, error) {

	return u.userRepo.CreateUser(ctx, user)
}

func (u *UseCase) UpdateUser(ctx context.Context, user models.User) error {
	return u.userRepo.UpdateUser(ctx, user)
}

func (u *UseCase) GetUser(ctx context.Context, id string) (*models.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.Wrap(err, "get user")
	}
	return u.userRepo.GetUser(ctx, id)
}

func (u *UseCase) DeleteUser(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.Wrap(err, "delete user")
	}
	return u.userRepo.DeleteUser(ctx, id)

}
