package domain

import (
	"time"

	"github.com/google/uuid"
)

type EmailLog struct {
	EventID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"event_id"`
	Subject   string    `json:"subject"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Body      string    `json:"body" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
