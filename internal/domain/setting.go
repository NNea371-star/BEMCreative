package domain

import (
	"time"
	"github.com/google/uuid"
)

type SiteSetting struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Key       string    `gorm:"uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WhatsappConfig struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PhoneNumber   string    `json:"phone_number"`
	OrderTemplate string    `gorm:"type:text" json:"order_template"`
	HireTemplate  string    `gorm:"type:text" json:"hire_template"`
	UpdatedAt     time.Time `json:"updated_at"`
}
