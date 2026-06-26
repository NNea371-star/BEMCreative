package domain

import (
	"time"
	"github.com/google/uuid"
)

type Portfolio struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	ClientName  string    `json:"client_name"`
	Year        int       `json:"year"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}