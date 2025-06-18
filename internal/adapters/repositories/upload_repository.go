package repositories

import (
	"context"
	"fmt"
	"os"

	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
	"gorm.io/gorm"
)

type uploadRepository struct {
	db *gorm.DB
}

func NewUploadRepository(db *gorm.DB) ports.UploadProfileRepository {
	return &uploadRepository{db: db}
}

func (r *uploadRepository) Save(ctx context.Context, upload *domain.UploadProfile) error {
	path := "internal/adapters/storage/uploads/profile_pictures/" + upload.ProfileImage
	err := os.WriteFile(path, upload.ImageData, 0644)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	return nil
}
