package repository

import (
	"BE/internal/domain"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateSession(s *domain.ChatSession) error {
	return r.db.Create(s).Error
}

func (r *ChatRepository) FindAllSessions() ([]domain.ChatSession, error) {
	var sessions []domain.ChatSession
	err := r.db.Order("created_at DESC").Find(&sessions).Error
	return sessions, err
}

func (r *ChatRepository) FindSessionByID(id string) (*domain.ChatSession, error) {
	var session domain.ChatSession
	err := r.db.Where("id = ?", id).First(&session).Error
	return &session, err
}

func (r *ChatRepository) FindMessagesBySessionID(sessionID string) ([]domain.ChatMessage, error) {
	var messages []domain.ChatMessage
	err := r.db.Where("session_id = ?", sessionID).Order("created_at ASC").Find(&messages).Error
	return messages, err
}

func (r *ChatRepository) CreateMessage(m *domain.ChatMessage) error {
	return r.db.Create(m).Error
}