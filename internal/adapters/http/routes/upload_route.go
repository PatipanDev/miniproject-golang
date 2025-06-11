package route

import (
	"github.com/PatipanDev/mini-project-golang/pkg/configs"

	"github.com/PatipanDev/mini-project-golang/internal/adapters/http/handlers"

	"github.com/PatipanDev/mini-project-golang/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func UploadRoutes(router fiber.Router, userhandler *handlers.HttpUserHandler) {
	user := router.Group("/uploads")

	user.Post("/profile/:id", userhandler.UploadProfilePicture)

	user.Use(middleware.JWTMiddleware(configs.SECRET_KEY))
}
