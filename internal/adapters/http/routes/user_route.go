package route

import (
	"os"
	"test-backend/internal/adapters/http/handlers"
	"test-backend/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router, userhandler *handlers.HttpUserHandler, authHandler *handlers.AuthHandler) {
	SECRET_KEY := os.Getenv("SECRET_KEY")
	user := router.Group("/users")

	user.Post("/register", userhandler.RegisterUser)
	user.Post("/login", authHandler.Login)
	//user.Get("/", userhandler.SearchUsers)
	user.Get("/", userhandler.FindUsers)
	user.Get("/pagination", userhandler.GetPaginationUsers)
	user.Get("/search", userhandler.GetUsers)

	user.Use(middleware.JWTMiddleware(SECRET_KEY))
}
