package handler

import (
	"errors"

	"github.com/P04KA/API/internal/apperr"
	"github.com/P04KA/API/internal/grpc/client"
	"github.com/P04KA/API/internal/models"
	"github.com/P04KA/API/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handle struct {
	uc usecase.UserProvider
}

var validate = validator.New()

func New(uc usecase.UserProvider) *Handle {
	return &Handle{uc: uc}
}
func (h *Handle) GetHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := validate.Var(id, "uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid format"})
	}

	user, err := h.uc.GetUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user"})
		}
		if errors.Is(err, apperr.ErrNotFound) {
			return apperr.ErrNotFound
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}
func (h *Handle) PostHandler(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := h.uc.CreateUser(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Can't save user"})
	}

	return c.JSON(fiber.Map{
		"message": "User has been created",
		"user":    id,
	})
}

func (h *Handle) DeleteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.uc.DeleteUser(c.Context(), id)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "delete user"})
		}
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

	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.uc.UpdateUser(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

	}
	return c.JSON(fiber.Map{
		"message": "User has been updated ",
		"user":    user,
	})
}

// grpc
func (h *Handle) GetUserStats(c *fiber.Ctx) error {
	period := c.Query("period", "day")

	resp, err := client.GetStats(c.Context(), period)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"users_created": resp.UsrCreated,
		"users_updated": resp.UsrUpdated,
		"users_deleted": resp.UsrDeleted,
		"period":        resp.Period,
	})
}
