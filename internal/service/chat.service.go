package service

import (
	"BE/internal/domain"
	"BE/internal/hub"
	"BE/internal/repository"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ChatService struct {
	repo *repository.ChatRepository
}

func NewChatService(repo *repository.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) CreateSession(visitorName, visitorWA string) (*domain.ChatSession, error) {
	if visitorName == "" {
		return nil, errors.New("nama pengunjung wajib diisi")
	}

	session := &domain.ChatSession{
		ID:          uuid.New(),
		VisitorName: visitorName,
		VisitorWA:   visitorWA,
		Status:      "active",
	}

	if err := s.repo.CreateSession(session); err != nil {
		return nil, err
	}

	// Notif ke admin via WebSocket
	hub.H.BroadcastToAdmins(fiber.Map{
		"type":    "chat:session_new",
		"message": "Pengunjung baru: " + session.VisitorName,
		"session": session,
	})

	return session, nil
}

func (s *ChatService) GetAllSessions() ([]domain.ChatSession, error) {
	return s.repo.FindAllSessions()
}

func (s *ChatService) GetSessionMessages(sessionID string) ([]domain.ChatMessage, error) {
	_, err := s.repo.FindSessionByID(sessionID)
	if err != nil {
		return nil, errors.New("sesi chat tidak ditemukan")
	}
	return s.repo.FindMessagesBySessionID(sessionID)
}

func (s *ChatService) SaveMessage(sessionID, senderRole, message string) (*domain.ChatMessage, error) {
	msg := &domain.ChatMessage{
		ID:         uuid.New(),
		SessionID:  uuid.MustParse(sessionID),
		SenderRole: senderRole,
		Message:    message,
	}
	return msg, s.repo.CreateMessage(msg)
}