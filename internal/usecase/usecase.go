package usecase

import (
	"context"
	"errors"

	"github.com/P04KA/API.git/internal/models"
	"github.com/P04KA/API.git/internal/repository"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type UseCase struct {
	userRepo *repository.UserRepo
	validate *validator.Validate
}

func New(userRepo *repository.UserRepo) *UseCase {
	return &UseCase{
		userRepo: userRepo,
		validate: validator.New(),
	}
}

func (u *UseCase) CreateUser(ctx context.Context, user models.User) (string, error) {
	if err := u.validate.Struct(user); err != nil {
		return "", errors.New("validate fail")
	}
	return u.userRepo.CreateUser(ctx, user)
}

func (u *UseCase) UpdateUser(ctx context.Context, user models.User) error {
	if err := u.validate.Struct(user); err != nil {
		return errors.New("validate fail")
	}
	if _, err := uuid.Parse(user.ID); err != nil {
		return errors.New("invalid user")
	}
	err := u.userRepo.UpdateUser(ctx, user)

	if err != nil && err.Error() == "user not found" {
		return errors.New("user not found") // Возвращаем ту же ошибку, но гарантируем ее наличие
	}
	return err
}

func (u *UseCase) GetUser(ctx context.Context, id string) (*models.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid user")
	}
	user, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, errors.New("user not found") // Преобразуем ошибку БД в бизнес-ошибку
	}
	return user, nil
}

func (u *UseCase) DeleteUser(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid user")
	}

	err := u.userRepo.DeleteUser(ctx, id)
	if err != nil && err.Error() == "user not found" {
		return errors.New("user not found")
	}
	return err
}
