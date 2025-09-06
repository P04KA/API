package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Age     string `json:"age"`
	Country string `json:"country"`
}

var users = map[string]User{}

func handler(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func getHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	user, ok := users[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Пользователя нет"})
	}
	return c.JSON(user)
}

func postHandler(c *fiber.Ctx) error {
	var user User
	// Сделать валидацию на age
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user.ID = uuid.New().String()
	users[user.ID] = user

	return c.JSON(fiber.Map{
		"message": "Пользователь cоздан",
		"user":    user,
	})
}

func deleteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := users[id]; !err {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Пользователя не существует"})
	}
	delete(users, id)
	return c.JSON(fiber.Map{
		"message": "Пользователь удален",
		"id":      id,
	})
}
func putHandler(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	users[user.ID] = user

	return c.JSON(fiber.Map{
		"message": "Пользователь обновлен",
		"user":    user,
	})
}

func main() {
	app := fiber.New()

	app.Get("/", handler)

	app.Get("/user/:id", getHandler)

	app.Post("/user", postHandler)

	app.Put("/user", putHandler)

	app.Delete("/user/:id", deleteHandler)
	if err := app.Listen(":3000"); err != nil {
		os.Exit(1)
	}

}
