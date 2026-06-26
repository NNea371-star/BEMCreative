package domain

import (
	"time"
	"github.com/google/uuid"
)

type ProductCategory struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CategoryID  uuid.UUID       `gorm:"type:uuid;not null" json:"category_id"`
	Category    ProductCategory `gorm:"foreignKey:CategoryID" json:"category"`
	Name        string          `gorm:"not null" json:"name"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Stock       int             `gorm:"default:0" json:"stock"`
	ImageURL    string          `json:"image_url"`
	IsAvailable bool            `gorm:"default:true" json:"is_available"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}