package main

import (
	"github.com/gofiber/fiber/v2"
)

/*
	func main() {
		app := fiber.New(fiber.Config{
			Prefork:       true,    // буст
			ServerHeader:  "Fiber", // заголовок
			CaseSensitive: true,    // чек URL
			StrictRouting: true,    // маршруты
		})
		type User struct {
			Name    string `json:"name"`
			Age     string `json:"age"`
			Country string `json:"country"`
		}
		app.Post("/user", func(c *fiber.Ctx) error {
			var user User
			if errors := c.BodyParser(&user); errors != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.Error()})
			}
			return c.JSON(fiber.Map{
				"message": "Пользователь получен",
				"user":    user,
			})
		})
		app.Get("/user/:id", func(c *fiber.Ctx) error {
			//	country := "Russia"
			//	name := "Tema"
			//	age := "19"
			id := c.Params("id")
			return c.SendString("ID: " + id)
		})
		app.Delete("/user:id", func(c *fiber.Ctx) error {
			delete(fiber.Map)
		})

		app.Listen(":3000")
	}
*/

type User struct {
	Name    string `json:"name"`
	Age     string `json:"age"`
	Country string `json:"country"`
}

var handler = func(c *fiber.Ctx) error { return nil }
var Get_handler = func(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.SendString("ID: " + id)
}
var Post_handler = func(c *fiber.Ctx) error {
	var user User
	if errors := c.BodyParser(&user); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.Error()})
	}
	return c.JSON(fiber.Map{
		"message": "Пользователь получен",
		"user":    user,
	})
}

func main() {
	app := fiber.New()

	app.Get("/", handler)

	app.Get("/user/:id", Get_handler)

	app.Post("/user", Post_handler)

	app.Put("/user", handler)

	app.Delete("/user:id", handler)

	app.Listen(":3000")

}
