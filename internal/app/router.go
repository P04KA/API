package app

import (
	"github.com/P04KA/API/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func getRouter(handle *handler.Handle) *fiber.App {
	app := fiber.New()

	app.Get("/user/:id", handle.GetHandler)

	app.Post("/user", handle.PostHandler)

	app.Put("/user", handle.PutHandler)

	app.Delete("/user/:id", handle.DeleteHandler)
	return app
}
