package repositories

import (
	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
	"gorm.io/gorm"
)

// adapters/repositories/email_log_repository.go
type EmailLogRepository struct {
	db *gorm.DB
}

func NewEmailLogRepository(db *gorm.DB) ports.EmailLogRepository {
	return &EmailLogRepository{db: db}
}

func (r *EmailLogRepository) Save(log *domain.EmailLog) error {
	return r.db.Create(log).Error
}
