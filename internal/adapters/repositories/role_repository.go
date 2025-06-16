// internal/adapters/repositories/role_repository.go
package repositories

import (
	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"gorm.io/gorm"
)

type GormRoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *GormRoleRepository {
	return &GormRoleRepository{db: db}
}

// 🔍 ค้นหา Role ตามชื่อ
func (r *GormRoleRepository) FindByName(name string, role *domain.Role) error {
	return r.db.Where("name = ?", name).First(role).Error
}

// ➕ สร้าง Role ใหม่
func (r *GormRoleRepository) Create(role *domain.Role) error {
	return r.db.Create(role).Error
}
