package handlers

import (
	"context"
	"database/sql"
	"io"
	"log"

	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

type UploadHandler struct {
	service     ports.UplaodProfileService
	db          *gorm.DB
	riverClient *river.Client[*sql.Tx]
}

func NewUploadHandler(service ports.UplaodProfileService, db *gorm.DB, riverClient *river.Client[*sql.Tx]) *UploadHandler {
	return &UploadHandler{service: service, db: db, riverClient: riverClient}
}

func (h *UploadHandler) UploadProfile(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required as query parameter or path parameter"})
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error parse strign to uuid"})
	}
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "image is required")
	}
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	upload := domain.UploadProfile{
		ID:           id,
		ProfileImage: fileHeader.Filename,
		ImageData:    imageData,
	}

	tx := h.db.Begin()
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
	sqlTx := tx.Statement.ConnPool.(*sql.Tx)

	_, err = h.riverClient.InsertTx(context.Background(), sqlTx, upload, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatal(err)
	}

	return c.JSON(fiber.Map{"status": "job enqueued"})
}
