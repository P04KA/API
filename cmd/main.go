package main

import (
	"os"

	"log/slog"

	"github.com/P04KA/API.git/database"
	"github.com/P04KA/API.git/internal/app"
)

func main() {

	if err := database.Migrate("postgresql://postgres:postgres@postgres:5432/postgres"); err != nil {
		slog.Error("migrate database", slog.Any("error", err))
		os.Exit(1)
	}
	if err := app.Run(); err != nil {

		slog.Error("run app", slog.Any("error", err))
		os.Exit(1)
	}
}
