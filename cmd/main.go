package main

import (
	"os"

	"log/slog"

	"github.com/P04KA/API.git/internal/app"
)

func main() {
	if err := app.Run(); err != nil {

		slog.Error("run app", slog.Any("error", err))
		os.Exit(1)
	}
}
