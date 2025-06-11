package repositories

import (
	"fmt"

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

func (r *GormUserRepository) Get() ([]domain.User, error) {
	var users []domain.User

	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
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
func (r *GormUserRepository) SearchData(query string) ([]domain.User, error) {
	var users []domain.User

	if err := r.db.Where("concat(first_name, ' ', last_name) ILIKE ?", "%"+query+"%").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// pagination
func (r *GormUserRepository) FindAll(offset int, limit int) ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (r *GormUserRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Count(&count).Error
	return count, err
}

// search
func (r *GormUserRepository) FindUsers(filter *domain.UserFilter) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	db := r.db.Model(&domain.User{})

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		db = db.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?", search, search, search)
	}

	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	db.Count(&total)

	offset := (filter.Page - 1) * filter.Limit
	err := db.Limit(filter.Limit).Offset(offset). /*.Order("created_at DESC")*/ Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *GormUserRepository) UpdateUserProfilePicURL(id string, url string) error {
	result := r.db.Model(&domain.User{}).Where("id = ?", id).Update("profile_image", url)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %s not found", id)
	}
	return nil
}
