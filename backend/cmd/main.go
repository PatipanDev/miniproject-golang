package main

import (
	"log"

	"github.com/PatipanDev/mini-project-golang/configs"
	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
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

	err = fiberServ.Listen(":" + configs.SERVER_HTTP_PORT)
	if err != nil {
		log.Fatal("Failed to start Fiber HTTP Server:", err)
	}

}
