package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrderLog untuk menyimpan pesanan dari WhatsApp
type OrderLog struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	VisitorName  string    `gorm:"not null" json:"visitor_name"`
	VisitorWA    string    `gorm:"not null" json:"visitor_wa"`
	ProductName  string    `json:"product_name"`
	ProjectType  string    `json:"project_type"`
	Budget       string    `json:"budget"`
	Description  string    `json:"description"`
	Status       string    `gorm:"default:'pending'" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (OrderLog) TableName() string {
	return "order_logs"
}

// Status constants
const (
	OrderStatusPending   = "pending"
	OrderStatusProcessed = "processed"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)