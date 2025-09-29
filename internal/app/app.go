package app

import (
	"github.com/P04KA/API/database"
	"github.com/P04KA/API/internal/cache"
	"github.com/P04KA/API/internal/handler"
	"github.com/P04KA/API/internal/repository"
	"github.com/P04KA/API/internal/storage"
	"github.com/P04KA/API/internal/usecase"
	"github.com/pkg/errors"
)

func Run() error {

	if err := database.Migrate("postgresql://postgres:postgres@postgres:5432/postgres"); err != nil {
		return errors.Wrap(err, "migrate")
	}

	conn, err := storage.GetConnect("postgresql://postgres:postgres@postgres:5432/postgres")
	if err != nil {
		return errors.Wrap(err, "Connect")
	}
	defer conn.Close()

	userRepo := repository.New(conn)
	cacheDecorator := cache.NewDecorator(userRepo)
	uc := usecase.New(cacheDecorator)
	handle := handler.New(uc)
	app := getRouter(handle)

	if err := app.Listen(":8080"); err != nil {
		return errors.Wrap(err, "start app")
	}
	return nil

}
