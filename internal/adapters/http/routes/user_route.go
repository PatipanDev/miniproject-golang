package route

import (
	"test-backend/internal/adapters/http/handlers"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router, userhandler *handlers.HttpUserHandler, authHandler *handlers.AuthHandler) {
	user := router.Group("/users")

	user.Post("/register", userhandler.RegisterUser)
	user.Post("/login", authHandler.Login)
}
