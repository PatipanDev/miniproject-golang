package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type USER_STATUS string
type USER_ROLE string

const (
	USER_STATUS_ACTIVE   USER_STATUS = "active"
	USER_STATUS_INACTIVE USER_STATUS = "inactive"
	USER_STATUS_BANNED   USER_STATUS = "banned"
)

const (
	USER_ROLE_ADMIN    USER_ROLE = "admin"
	USER_ROLE_PREPARER USER_ROLE = "preparer"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Email     string         `json:"email" validate:"required,max=50" gorm:"size:50"`
	Username  string         `json:"username" validate:"required,max=255" gorm:"size:255"`
	Password  string         `json:"password"`
	Status    USER_STATUS    `json:"status" validate:"required,max=50" gorm:"size:50"`
	Role      USER_ROLE      `json:"role" validate:"required,max=50" gorm:"size:50"`
}
