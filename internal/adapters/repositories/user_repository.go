package repositories

import (
	"context"
	"fmt"

	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db          *gorm.DB
	minioClient *minio.Client
	bucketName  string
}

func NewUserRepository(db *gorm.DB, minioClient *minio.Client, bucketName string) ports.UserRepository {
	return &GormUserRepository{
		db:          db,
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (r *GormUserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) Get() ([]domain.User, error) {
	var users []domain.User

	if err := r.db.Preload("Roles").Find(&users).Error; err != nil {
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
	if err := r.db.Preload("Roles").
		Select("id", "email", "first_name", "last_name", "employee_id", "status", "profile_image", "created_at", "updated_at", "deleted_at").
		First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

/*func (r *GormUserRepository) GetProfile(id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.
		First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}*/

func (r *GormUserRepository) SearchData(query string) ([]domain.User, error) {
	var users []domain.User

	if err := r.db.Preload("Roles").Where("concat(first_name, ' ', last_name) ILIKE ?", "%"+query+"%").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// pagination
func (r *GormUserRepository) FindAll(offset int, limit int) ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Preload("Roles").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
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
		db = db.Preload("Roles").Where("concat(first_name, ' ', last_name) LIKE ? OR employee_id LIKE ? OR email LIKE ? OR status LIKE ?", search, search, search, search)
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

func (r *GormUserRepository) UpdateUserProfilePicName(id string, filename string) error {
	var user domain.User

	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if user.ProfileImage != "" {
		oldObjectName := user.ProfileImage
		fmt.Println("Object", oldObjectName)
		err := r.minioClient.RemoveObject(context.Background(), r.bucketName, oldObjectName, minio.RemoveObjectOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete old profile image from MinIO: %w", err)
		}
	}

	updateResult := r.db.Model(&user).Update("profile_image", filename)
	if updateResult.Error != nil {
		return updateResult.Error
	}
	if updateResult.RowsAffected == 0 {
		return fmt.Errorf("user with ID %s not found", id)
	}
	return nil
}
