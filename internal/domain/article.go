package domain

import (
	"time"
	"github.com/google/uuid"
)

type BlogCategory struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type Article struct {
	ID           uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AuthorID     uuid.UUID    `gorm:"type:uuid;not null" json:"author_id"`
	Author       User         `gorm:"foreignKey:AuthorID" json:"author"`
	CategoryID   uuid.UUID    `gorm:"type:uuid;not null" json:"category_id"`
	Category     BlogCategory `gorm:"foreignKey:CategoryID" json:"category"`
	Title        string       `gorm:"not null" json:"title"`
	Slug         string       `gorm:"uniqueIndex;not null" json:"slug"`
	Content      string       `gorm:"type:text" json:"content"`
	ThumbnailURL string       `json:"thumbnail_url"`
	Status       string       `gorm:"default:draft" json:"status"`
	PublishedAt  *time.Time   `json:"published_at"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}