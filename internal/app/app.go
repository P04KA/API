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

	if err := app.Listen(":8080"); err != nil {
		return errors.Wrap(err, "start app")
	}
	return nil

}
