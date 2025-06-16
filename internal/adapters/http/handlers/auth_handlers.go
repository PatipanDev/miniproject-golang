package handlers

import (
	"fmt"
	"time"

	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
	"github.com/PatipanDev/mini-project-golang/pkg/configs"
	"github.com/golang-jwt/jwt/v5"

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

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})
	return c.JSON(fiber.Map{
		"message": "Login successful!",
	})
}

func AuthRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	jwtSecret := configs.SECRET_KEY

	token, err := jwt.ParseWithClaims(cookie /*,jwt.MapClaims{}*/, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	//claim := token.Claims.(jwt.MapClaims)

	//fmt.Println(claim)
	fmt.Println(token)

	return c.Next()
}
