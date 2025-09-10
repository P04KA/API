package main

import (
	"os"

	"log/slog"

	"github.com/P04KA/API.git/internal/storage"
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

func main() {
	app := fiber.New()
	// Глянуть про указатели
	conn, err := storage.GetConnect("postgresql://postgres:postgres@postgres:5432/postgres")
	if err != nil {
		slog.Error("db conn", slog.Any("error", err))
		os.Exit(1)
	}
	defer conn.Close()

	app.Get("/user/:id", getHandler)

	app.Post("/user", postHandler)

	app.Put("/user", putHandler)

	app.Delete("/user/:id", deleteHandler)
	if err := app.Listen(":3000"); err != nil {
		//log ошибки
		os.Exit(1)
	}

}
