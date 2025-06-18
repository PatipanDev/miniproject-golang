package domain

import (
	"fmt"
	"math/rand"
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
	USER_ROLE_EMPLOYEE USER_ROLE = "employee"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	EmployeeID string `json:"employee_id" gorm:"uniqueIndex;size:10"`

	FirstName string `json:"first_name" validate:"required,max=255" gorm:"size:255"`
	LastName  string `json:"last_name" validate:"required,max=255" gorm:"size:255"`
	//PhoneNumber string `jso:"phone" validate:"required,max=255" gorm:"size:255"`
	//Birthday time.Time `json:"birthday"`

	Email        string      `json:"email" validate:"required,max=50" gorm:"size:50"`
	Username     string      `json:"username" validate:"required,max=255" gorm:"size:255"`
	Password     string      `json:"password"`
	Status       USER_STATUS `json:"status" validate:"required,max=50" gorm:"size:50"`
	Roles        []Role      `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID;" json:"roles"`
	ProfileImage string      `json:"profile_image,omitempty" gorm:"type:text"`
	//Work Work `json:"work"`
}

type Role struct {
	ID   uint      `gorm:"primaryKey" json:"id"`
	Name USER_ROLE `gorm:"unique;size:50" json:"name"` // เช่น "admin", "preparer"
	//User []User    `gorm:"many2many:user_roles;"`
}

// create employee id
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// วนสุ่มจนกว่าจะได้ EmployeeID ที่ไม่ซ้ำใน DB
	for {
		code := generateRandomEmployeeCode()

		var count int64
		tx.Model(&User{}).Where("employee_id = ?", code).Count(&count)

		if count == 0 {
			u.EmployeeID = code
			break
		}
	}
	return nil
}

func generateRandomEmployeeCode() string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	digits := []rune("0123456789")
	rand.Seed(time.Now().UnixNano())

	// สุ่มตัวอักษร 2 ตัว
	first := letters[rand.Intn(len(letters))]
	second := letters[rand.Intn(len(letters))]

	// สุ่มตัวเลข 1 ตัว
	digit := digits[rand.Intn(len(digits))]

	return fmt.Sprintf("#%c%c%c", first, second, digit)
}

type UserFilter struct {
	Search string `json:"search"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Status string `json:"status"`
}

type ResUSerNoID struct {
	FullName   string    `json:"full_name"`
	EmployeeID string    `json:"employee_id"`
	Email      string    `json:"Email"`
	Status     string    `json:"status"`
	UpdatedAt  time.Time `json:"update_at"`
	Roles      []string  `json:"role"`
}
type ResUserWaithID struct {
	ID         uuid.UUID `json:"id"`
	FullName   string    `json:"full_name"`
	EmployeeID string    `json:"employee_id"`
	Email      string    `json:"Email"`
	Status     string    `json:"status"`
	UpdatedAt  time.Time `json:"update_at"`
	Roles      []string  `json:"role"`
}

type UserProfileResponse struct {
	ID           uuid.UUID      `json:"id"`
	ProfileImage string         `json:"profile_image"`
	FullName     string         `json:"full_name"`
	EmployeeID   string         `json:"employee_id"`
	Email        string         `json:"email"`
	Status       string         `json:"status"`
	Roles        []string       `json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

type UploadProfile struct {
	ID           uuid.UUID `json:"id"`
	ProfileImage string    `json:"profile_image"`
	ImageData    []byte    `json:"image_data"`
}

func (UploadProfile) Kind() string {
	return "upload_image"
}
