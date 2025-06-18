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

func (h *AuthHandler) CookieLogin(c *fiber.Ctx) error {
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
		Secure:   false,
		SameSite: "Lax",
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful!",
	})
}

func (h *AuthHandler) JwtLogin(c *fiber.Ctx) error {
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

func AuthRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	jwtSecret := configs.SECRET_KEY

	fmt.Println("Received JWT Cookie:", cookie)

	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing jwt cookie"})
	}

	token, err := jwt.ParseWithClaims(cookie /*,&jwt.RegisteredClaims{}*/, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claim := token.Claims.(jwt.MapClaims)
	userID := claim["sup"].(string)
	fullname := claim["fullname"].(string)

	fmt.Println("Token Claims:", claim)
	c.Locals("sup", userID)
	c.Locals("fullname", fullname)
	//fmt.Println(token)

	return c.Next()
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})
	return c.JSON(fiber.Map{"message": "Logged out"})
}

func GetProfile(c *fiber.Ctx) error {
	fullname := c.Locals("fullname").(string)
	return c.JSON(fiber.Map{"message": "Welcome!", "fullname": fullname})
}

func GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return c.JSON(fiber.Map{
		"user_id": userID,
	})
}
