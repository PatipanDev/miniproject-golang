package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

// adapters/repositories/email_log_repository.go
type EmailLogRepository struct {
	db          *gorm.DB
	riverClient *river.Client[*sql.Tx]
}

func NewEmailLogRepository(db *gorm.DB, riverClient *river.Client[*sql.Tx]) ports.EmailLogRepository {
	return &EmailLogRepository{db: db}
}

func (r *EmailLogRepository) Save(log *domain.EmailLog) error {
	return r.db.Create(log).Error
}

func (r *EmailLogRepository) ProcessUserEmailJob(args *domain.SendEmailArgs) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	sqlTx := tx.Statement.ConnPool.(*sql.Tx)

	if _, err := r.riverClient.InsertTx(context.Background(), sqlTx, args, nil); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
