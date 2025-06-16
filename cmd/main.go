package main

import (
	"log"

	"github.com/PatipanDev/mini-project-golang/internal/adapters/http/handlers"
	routers "github.com/PatipanDev/mini-project-golang/internal/adapters/http/routes"
	"github.com/PatipanDev/mini-project-golang/internal/adapters/repositories"
	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/services"
	"github.com/PatipanDev/mini-project-golang/pkg/configs"
	"github.com/PatipanDev/mini-project-golang/pkg/smtp"
)

func main() {

	param := configs.NewFiberHttpServiceParams()
	fiberServ := configs.NewFiberHTTPService(param)

	db, err := configs.NewDatabase()
	if err != nil {
		log.Fatal("Failed to start Database:", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(&domain.User{}, &domain.Role{}, &domain.EmailLog{})
	if err != nil {
		log.Fatal("Failed to auto migrate models:", err)
	}

	userRepo := repositories.NewUerRepository(db)
	fileRepo := repositories.NewFileStorageRepository("internal/adapters/storage", configs.BASE_URL)

	// EmailLogRepository & EmailSender
	emailLogRepo := repositories.NewEmailLogRepository(db)
	emailSender := smtp.NewSMTPSender(emailLogRepo)

	roleRepo := repositories.NewRoleRepository(db)

	userServ := services.NewUserService(userRepo, fileRepo, emailSender, roleRepo)
	userHandler := handlers.NewHttpUserHandler(userServ)

	authServ := services.NewAurhService(userRepo, configs.SECRET_KEY)
	authHandler := handlers.NewAuthHandler(authServ)

	api := fiberServ.Group("/api")
	routers.UserRoutes(api, userHandler, authHandler)
	routers.UploadRoutes(api, userHandler)

	fiberServ.Static("uploads", "internal/adapters/storage/uploads")

	err = fiberServ.Listen(":" + configs.SERVER_HTTP_PORT)
	if err != nil {
		log.Fatal("Failed to start Fiber HTTP server:", err)
	}
}
