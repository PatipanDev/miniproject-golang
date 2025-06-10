package domain

import (
	"time"
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
	USER_ROLE_PREPARER USER_ROLE = "perparer"
	USER_ROLE_USER     USER_ROLE = "user"
)

type User struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  *time.Time  `json:"deleted_at,omitempty"`
	FirstName  string      `json:"first_name" validate:"required,max=255" gorm:"size:255"`
	LastName   string      `json:"last_name" validate:"required,max=255"`
	Email      string      `json:"email" validate:"required,max=150" gorm:"size:50"`
	Username   string      `json:"username" validate:"required,max=255" gorm:"size:255"`
	Password   string      `json:"password"`
	Role       USER_ROLE   `json:"role" validate:"required,max=50" gorm:"size:50"`
	Status     USER_STATUS `json:"status" validate:"required,max=50" gorm:"size:50"`
	EmployeeID string      `json:"employee_id"  validate:"required,max=255" gorm:"size:255"`
}

var TNUser = "users"

func (*User) TableName() string {
	return TNUser
}

type UserFilter struct {
	Search string `json:"search"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Status string `json:"status"`
}
