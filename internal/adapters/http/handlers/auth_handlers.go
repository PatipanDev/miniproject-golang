package handlers

import (
	"fmt"
	"test-backend/internal/core/domain"
	"test-backend/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req = &domain.User{}
	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"erorr": err.Error()})
	}
	return c.JSON(fiber.Map{
		"access_token": token,
	})
}
