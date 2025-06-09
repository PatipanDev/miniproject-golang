package repositories

import (
	"test-backend/internal/core/domain"
	"test-backend/internal/core/ports"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewUerRepository(db *gorm.DB) ports.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
