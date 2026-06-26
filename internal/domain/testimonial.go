package domain

import (
	"time"
	"github.com/google/uuid"
)

type Testimonial struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ClientName    string    `gorm:"not null" json:"client_name"`
	ClientCompany string    `json:"client_company"`
	Message       string    `gorm:"type:text;not null" json:"message"`
	Rating        int       `gorm:"default:5" json:"rating"`
	IsVisible     bool      `gorm:"default:false" json:"is_visible"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}