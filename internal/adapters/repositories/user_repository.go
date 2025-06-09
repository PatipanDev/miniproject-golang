package repositories

import (
	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
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

func (r *GormUserRepository) Update(user *domain.User, id string) error {
	var existingUser domain.User
	if err := r.db.First(&existingUser, "id = ?", id).Error; err != nil {
		return err
	}

	if err := r.db.Model(&existingUser).Updates(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormUserRepository) Delete(id string) error {
	var user domain.User
	if err := r.db.First(&user, "id =?", id).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormUserRepository) FindUserById(id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.
		Select("id", "email", "username", "status", "role", "profile_image", "created_at", "updated_at", "deleted_at").
		First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
