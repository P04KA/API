package app

import (
	"time"

	"github.com/P04KA/API/config"
	"github.com/P04KA/API/database"
	"github.com/P04KA/API/internal/cache"
	"github.com/P04KA/API/internal/handler"
	"github.com/P04KA/API/internal/repository"
	"github.com/P04KA/API/internal/storage"
	"github.com/P04KA/API/internal/usecase"
	"github.com/gofiber/contrib/circuitbreaker"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func Run() error {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return errors.Wrap(err, "cfg load")
	}

	if err := database.Migrate(cfg.DB.DBURL); err != nil {
		return errors.Wrap(err, "migrate")
	}

	conn, err := storage.GetConnect(cfg.DB.DBURL)
	if err != nil {
		return errors.Wrap(err, "Connect")
	}
	defer conn.Close()

	userRepo := repository.New(conn)
	cacheDecorator := cache.NewDecorator(10*time.Minute, 1*time.Minute, userRepo)
	uc := usecase.New(cacheDecorator)
	handle := handler.New(uc)
	app := getRouter(handle)

	cb := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 3,
		Timeout:          10 * time.Second,
		OnOpen: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusServiceUnavailable).
				JSON(fiber.Map{"error": "Circuit Open: Service unavailable"})
		},
		OnHalfOpen: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).
				JSON(fiber.Map{"error": "Circuit Half-Open: Retrying service"})
		},
		OnClose: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).
				JSON(fiber.Map{"message": "Circuit Closed: Service recovered"})
		},
	})

	app.Get("/user/:id", circuitbreaker.Middleware(cb), func(c *fiber.Ctx) error {
		return c.SendString("Post and Update services is protected by a Circuit Breaker")
	})
	app.Get("/user", circuitbreaker.Middleware(cb), func(c *fiber.Ctx) error {
		return c.SendString("Get and Delete is protected by a Circuit Breaker")
	})

	if err := app.Listen(":8082"); err != nil {
		return errors.Wrap(err, "start app")
	}
	return nil

}
