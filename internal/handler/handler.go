package handler

import (
	"github.com/P04KA/API/internal/models"
	"github.com/P04KA/API/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type Handle struct {
	uc *usecase.UseCase
}

func New(uc *usecase.UseCase) *Handle {
	return &Handle{uc: uc}
}
func (h *Handle) GetHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.uc.GetUser(c.Context(), id)
	if err != nil {
		if err.Error() == "invalid user" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user"})
		}
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
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

	id, err := h.uc.CreateUser(c.Context(), user)
	if err != nil {
		if err.Error() == "validate fail" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Can't save user"})
	}

	user.ID = id
	return c.JSON(fiber.Map{
		"message": "User has been created",
		"user":    user,
	})
}

func (h *Handle) DeleteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.uc.DeleteUser(c.Context(), id)

	if err != nil {
		if err.Error() == "invalid user" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user"})
		}

		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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

	if err := h.uc.UpdateUser(c.Context(), user); err != nil {
		if err.Error() == "validate fail" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err.Error() == "invalid user" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "User has been updated ",
		"user":    user,
	})
}
