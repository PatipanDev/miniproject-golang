package route

import (
	"github.com/PatipanDev/mini-project-golang/configs"

	"github.com/PatipanDev/mini-project-golang/internal/adapters/http/handlers"

	"github.com/PatipanDev/mini-project-golang/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router, userhandler *handlers.HttpUserHandler, authHandler *handlers.AuthHandler) {
	user := router.Group("/users")

	user.Post("/register", userhandler.RegisterUser)
	user.Post("/login", authHandler.Login)

	user.Put("/update/:id", userhandler.UpdateUser)
	user.Get("/getbyid/:id", userhandler.GetUserById)

	//user.Get("/", userhandler.SearchUsers)
	user.Get("/", userhandler.FindUsers)
	user.Get("/pagination", userhandler.GetPaginationUsers)
	user.Get("/search", userhandler.GetUsers)
	user.Delete("/delete/:id", userhandler.DeleteUser)

	user.Use(middleware.JWTMiddleware(configs.SECRET_KEY))
}
