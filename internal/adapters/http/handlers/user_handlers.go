package handlers

import (
	"test-backend/internal/core/domain"
	"test-backend/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type HttpUserHandler struct {
	userservice ports.UserService
}

func NewHttpUserHandler(userservice ports.UserService) *HttpUserHandler {
	return &HttpUserHandler{userservice: userservice}
}

func (h *HttpUserHandler) RegisterUser(c *fiber.Ctx) error {
	user := new(domain.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.userservice.RegisterUser(user); err != nil {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}
