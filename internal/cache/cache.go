package cache

import (
	"context"
	"sync"
	"time"

	"github.com/P04KA/API/internal/models"
	"github.com/P04KA/API/internal/repository"
	"github.com/pkg/errors"
)

type CacheDecorator struct {
	userRepo   repository.UserProvider
	user       map[string]*cacheItem
	mu         sync.RWMutex
	defaultTTL time.Duration
}

type cacheItem struct {
	user      *models.User
	updatedAt time.Time
}

func NewDecorator(defaultTTL, cleanupInterval time.Duration, repo repository.UserProvider) *CacheDecorator {
	c := &CacheDecorator{
		userRepo:   repo,
		user:       make(map[string]*cacheItem),
		defaultTTL: defaultTTL,
	}

	go c.Cleanup(cleanupInterval)

	return c
}

func (c *CacheDecorator) Cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.user {
			if now.After(item.updatedAt.Add(c.defaultTTL)) {
				delete(c.user, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *CacheDecorator) GetUser(ctx context.Context, id string) (*models.User, error) {
	c.mu.RLock()
	item, ok := c.user[id]
	c.mu.RUnlock()

	if ok && time.Now().Before(item.updatedAt.Add(c.defaultTTL)) {
		return item.user, nil
	}

	user, err := c.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "get user from db")
	}

	c.mu.Lock()
	c.user[id] = &cacheItem{
		user:      user,
		updatedAt: time.Now(),
	}
	c.mu.Unlock()
	return user, nil
}

func (c *CacheDecorator) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	createdUser, err := c.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "create user in db")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.user[createdUser.ID] = &cacheItem{
		user:      createdUser,
		updatedAt: time.Now(),
	}
	return createdUser, nil
}

func (c *CacheDecorator) UpdateUser(ctx context.Context, user models.User) error {
	err := c.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, "update user in db")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.user[user.ID] = &cacheItem{
		user:      &user,
		updatedAt: time.Now(),
	}
	return nil
}

func (c *CacheDecorator) DeleteUser(ctx context.Context, id string) error {
	err := c.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return errors.Wrap(err, "delete user from db")
	}

	c.mu.Lock()

	delete(c.user, id)
	return nil
}

func (c *CacheDecorator) CountUserCreated(ctx context.Context, time time.Time) (int64, error) {
	return c.userRepo.CountUserCreated(ctx, time)
}

func (c *CacheDecorator) CountUserUpdated(ctx context.Context, time time.Time) (int64, error) {
	return c.userRepo.CountUserUpdated(ctx, time)
}

func (c *CacheDecorator) CountUserDeleted(ctx context.Context, time time.Time) (int64, error) {
	return c.userRepo.CountUserDeleted(ctx, time)
}
