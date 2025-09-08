package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Age     string `json:"age"`
	Country string `json:"country"`
}

var conn *pgxpool.Pool

func handler(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func getHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	var user User
	err := conn.QueryRow(c.Context(), "SELECT id, name, age, country FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Age, &user.Country)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
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
	_, err := conn.Exec(c.Context(), "INSERT INTO users (id, name, age, country) VALUES ($1, $2, $3, $4)", user.ID, user.Name, user.Age, user.Country)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Can't save user"})
	}

	return c.JSON(fiber.Map{
		"message": "User has been created",
		"user":    user,
	})
}

func deleteHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := conn.Exec(c.Context(), "DELETE FROM users WHERE id = $1", id)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"message": "User has been removed",
		"id":      id,
	})
}
func putHandler(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := conn.Exec(c.Context(), "UPDATE users SET name = $2, age = $3, country = $4 WHERE id = $1", user.ID, user.Name, user.Age, user.Country)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "User has been updated ",
		"user":    user,
	})
}

func GetConnect() {
	var err error
	conn, err = pgxpool.New(context.Background(), "postgresql://user1:password@172.18.0.2:5432/db")
	if err != nil {
		println("Can't connect to DB", err.Error())
		os.Exit(1)
	}
}

func main() {
	app := fiber.New()

	GetConnect()
	defer conn.Close()

	app.Get("/", handler)

	app.Get("/user/:id", getHandler)

	app.Post("/user", postHandler)

	app.Put("/user", putHandler)

	app.Delete("/user/:id", deleteHandler)
	if err := app.Listen(":3000"); err != nil {
		os.Exit(1)
	}

}
