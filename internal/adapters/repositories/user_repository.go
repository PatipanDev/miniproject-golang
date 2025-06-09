package repositories

import (
	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
	"gorm.io/gorm"
)

type userRepositoryGorm struct {
	db *gorm.DB
}

func NewUserRepositoryGrom(db *gorm.DB) ports.UserRepository {
	return &userRepositoryGorm{db: db}
}

func (r *userRepositoryGorm) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return &user, err
	}

	return &user, nil
}

func (r *userRepositoryGorm) Create(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
