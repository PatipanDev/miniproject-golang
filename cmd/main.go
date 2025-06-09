package main

import (
	"log"
	"os"
	db "test-backend/DB"
	"test-backend/internal/adapters/http/handlers"
	route "test-backend/internal/adapters/http/routes"
	"test-backend/internal/adapters/repositories"
	"test-backend/internal/core/domain"
	"test-backend/internal/core/services"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	secret_key := os.Getenv("SECRET_KEY")

	app := fiber.New()
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Failed to start database")
	}

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatal("Migration error:", err)
	}

	userRepo := repositories.NewUerRepository(db)
	userServ := services.NewUserService(userRepo)
	userHandler := handlers.NewHttpUserHandler(userServ)

	authServ := services.NewAurhService(userRepo, secret_key)
	authHandler := handlers.NewAuthHandler(authServ)

	api := app.Group("/api")

	route.UserRoutes(api, userHandler, authHandler)

	err = app.Listen(":" + port)
	if err != nil {
		log.Fatal("Failed to start Fiber HTTP server:", err)
	}
}
