package handler

import (
	"errors"

	"github.com/P04KA/API.git/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handle struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Handle {
	return &Handle{conn: conn}
}

func (h *Handle) GetHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	var user models.User
	err := h.conn.QueryRow(c.Context(), "SELECT id, name, age, country FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Age, &user.Country)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}
func (h *Handle) PostHandler(c *fiber.Ctx) error {
	var user models.User
	// Сделать валидацию на age
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user.ID = uuid.New().String()
	_, err := h.conn.Exec(c.Context(), "INSERT INTO users (id, name, age, country) VALUES ($1, $2, $3, $4)", user.ID, user.Name, user.Age, user.Country)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Can't save user"})
	}

	return c.JSON(fiber.Map{
		"message": "User has been created",
		"user":    user,
	})
}
func (h *Handle) DeleteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	//валидацию на id

	_, err := h.conn.Exec(c.Context(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(fiber.Map{
		"message": "User has been removed",
		"id":      id,
	})
}
func (h *Handle) PutHandler(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	rows, err := h.conn.Exec(c.Context(), "UPDATE users SET name = $2, age = $3, country = $4 WHERE id = $1", user.ID, user.Name, user.Age, user.Country)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if rows.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(fiber.Map{
		"message": "User has been updated ",
		"user":    user,
	})
}
