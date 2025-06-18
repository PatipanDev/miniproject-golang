package route

import (
	"github.com/PatipanDev/mini-project-golang/pkg/configs"

	"github.com/PatipanDev/mini-project-golang/internal/adapters/http/handlers"

	"github.com/PatipanDev/mini-project-golang/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func UploadProfileRoutes(router fiber.Router, userhandler *handlers.UploadHandler) {
	upload := router.Group("/river")

	upload.Post("/profile/:id", userhandler.UploadProfile)

	upload.Use(middleware.JWTMiddleware(configs.SECRET_KEY))
}
