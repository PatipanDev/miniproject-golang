package main

import (
	"log"

	"github.com/PatipanDev/mini-project-golang/configs"
	"github.com/PatipanDev/mini-project-golang/internal/adapters/http/handlers"
	routers "github.com/PatipanDev/mini-project-golang/internal/adapters/http/routes"
	"github.com/PatipanDev/mini-project-golang/internal/adapters/repositories"
	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/services"
)

func main() {
	param := configs.NewFiberHttpServiceParams()
	fiberServ := configs.NewFiberHTTPService(param)

	db, err := configs.NewDatabase()
	if err != nil {
		log.Fatal("Failled to start Database:", err)

	}

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatal("Failed to connect to Database", err)
	}

	userRepo := repositories.NewUerRepository(db)
	userServ := services.NewUserService(userRepo)
	userHandler := handlers.NewHttpUserHandler(userServ)

	authServ := services.NewAurhService(userRepo, configs.SECRET_KEY)
	authHandler := handlers.NewAuthHandler(authServ)

	api := fiberServ.Group("/api")

	routers.UserRoutes(api, userHandler, authHandler)

	err = fiberServ.Listen(":" + configs.SERVER_HTTP_PORT)
	if err != nil {
		log.Fatal("Failed to start Fiber HTTP server:", err)
	}
}
