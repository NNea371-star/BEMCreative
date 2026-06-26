package domain

import (
	"time"
	"github.com/google/uuid"
)

type ChatSession struct {
	ID          uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	VisitorName string        `gorm:"not null" json:"visitor_name"`
	VisitorWA   string        `json:"visitor_wa"`
	Status      string        `gorm:"default:active" json:"status"`
	Messages    []ChatMessage `gorm:"foreignKey:SessionID" json:"messages,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type ChatMessage struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SessionID  uuid.UUID `gorm:"type:uuid;not null" json:"session_id"`
	SenderRole string    `gorm:"not null" json:"sender_role"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}