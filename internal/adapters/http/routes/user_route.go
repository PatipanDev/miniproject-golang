package route

import (
	//"github.com/PatipanDev/mini-project-golang/pkg/configs"

	"github.com/PatipanDev/mini-project-golang/internal/adapters/http/handlers"

	//"github.com/PatipanDev/mini-project-golang/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router, userhandler *handlers.HttpUserHandler, authHandler *handlers.AuthHandler) {
	user := router.Group("/users")

	router.Post("/register", userhandler.RegisterUser)
	router.Post("/login", authHandler.CookieLogin)
	router.Post("/loginv2", authHandler.JwtLogin)
	user.Use(handlers.AuthRequired)
	//user.Use(middleware.JWTMiddleware(configs.SECRET_KEY))
	user.Put("/update/:id", userhandler.UpdateUser)
	user.Get("/getbyid/:id", userhandler.GetUserById)
	user.Get("/profile", handlers.GetProfile)
	user.Get("/me", handlers.GetMe)

	//user.Get("/", userhandler.SearchUsers)
	user.Get("/", userhandler.FindUsers)
	user.Get("/pagination", userhandler.GetPaginationUsers)
	user.Get("/search", userhandler.GetUsers)
	user.Delete("/delete/:id", userhandler.DeleteUser)
}
