package cache

import (
	"context"

	"github.com/P04KA/API/internal/models"
	"github.com/P04KA/API/internal/repository"
	"github.com/pkg/errors"
)

type CacheDecorator struct {
	userRepo repository.UserProvider
	user     map[string]*models.User
}

func NewDecorator(repo repository.UserProvider) *CacheDecorator {
	return &CacheDecorator{
		userRepo: repo,
		user:     make(map[string]*models.User),
	}
}

// методы такие же как в repo
func (c *CacheDecorator) GetUser(ctx context.Context, id string) (*models.User, error) {

	if user, ok := c.user[id]; ok {
		return user, nil
	}

	user, err := c.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "get user from db")
	}

	c.user[id] = user
	return user, nil
}

func (c *CacheDecorator) CreateUser(ctx context.Context, user models.User) (string, error) {
	{
		return c.userRepo.CreateUser(ctx, user)
	}
}

func (c *CacheDecorator) UpdateUser(ctx context.Context, user models.User) error {
	{
		return c.userRepo.UpdateUser(ctx, user)
	}
}

func (c *CacheDecorator) DeleteUser(ctx context.Context, id string) error {
	{
		return c.userRepo.DeleteUser(ctx, id)
	}
}
