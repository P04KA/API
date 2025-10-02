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
	expiresAt time.Time // когда запись устареет
}

func NewDecorator(defaultTTL, cleanupInterval time.Duration, repo repository.UserProvider) *CacheDecorator {
	c := &CacheDecorator{
		userRepo:   repo,
		user:       make(map[string]*cacheItem),
		defaultTTL: defaultTTL,
	}

	// Запускаем очистку с тикером
	go c.cleanupWithTicker(cleanupInterval)

	return c
}

// Очистка с тикером
func (c *CacheDecorator) cleanupWithTicker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C // ждем тик
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.user {
			if now.After(item.expiresAt) {
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

	// Если есть в кэше и не устарел
	if ok && time.Now().Before(item.expiresAt) {
		return item.user, nil
	}

	// Если нет в кэше или устарел - берем из базы
	user, err := c.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "get user from db")
	}

	// Сохраняем в кэш
	c.mu.Lock()
	c.user[id] = &cacheItem{
		user:      user,
		expiresAt: time.Now().Add(c.defaultTTL),
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
	c.user[createdUser.ID] = &cacheItem{
		user:      createdUser,
		expiresAt: time.Now().Add(c.defaultTTL),
	}
	c.mu.Unlock()

	return createdUser, nil
}

func (c *CacheDecorator) UpdateUser(ctx context.Context, user models.User) error {
	err := c.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, "update user in db")
	}

	c.mu.Lock()
	c.user[user.ID] = &cacheItem{
		user:      &user,
		expiresAt: time.Now().Add(c.defaultTTL),
	}
	c.mu.Unlock()

	return nil
}

func (c *CacheDecorator) DeleteUser(ctx context.Context, id string) error {
	err := c.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return errors.Wrap(err, "delete user from db")
	}

	c.mu.Lock()
	delete(c.user, id)
	c.mu.Unlock()

	return nil
}
