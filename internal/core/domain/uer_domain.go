package domain

import (
	"time"
)

type USER_STATUS string

const (
	USER_STATUS_ACTIVE   USER_STATUS = "active"
	USER_STATUS_INACTIVE USER_STATUS = "inactive"
	USER_STATUS_BANNED   USER_STATUS = "banned"
)

type User struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time  `json:"deleted_at,omitempty"`
	FirstName string      `json:"first_name" validate:"required,max=255" gorm:"size:255"`
	LastName  string      `json:"last_name" validate:"required,max=255"`
	Email     string      `json:"email" validate:"required,max=150" gorm:"size:50"`
	Username  string      `json:"username" validate:"required,max=255" gorm:"size:255"`
	Password  string      `json:"password"`
	Role      string      `json:"role" gorm:"default:user"`
	Status    USER_STATUS `json:"status" validate:"required,max=50" gorm:"size:50"`
}

var TNUser = "users"

func (*User) TableName() string {
	return TNUser
}
