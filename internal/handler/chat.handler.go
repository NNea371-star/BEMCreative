package handler

import (
	"BE/internal/domain"
	"BE/internal/service"
	"BE/database"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	service *service.ChatService
}

func NewChatHandler(service *service.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

func (h *ChatHandler) CreateSession(c *fiber.Ctx) error {
	var body struct {
		VisitorName string `json:"visitor_name"`
		VisitorWA   string `json:"visitor_wa"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}

	session, err := h.service.CreateSession(body.VisitorName, body.VisitorWA)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(201).JSON(session)
}

func (h *ChatHandler) GetAllSessions(c *fiber.Ctx) error {
	sessions, err := h.service.GetAllSessions()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(sessions)
}

func (h *ChatHandler) GetSessionMessages(c *fiber.Ctx) error {
	messages, err := h.service.GetSessionMessages(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(messages)
}

func (h *ChatHandler) GetDashboardStats(c *fiber.Ctx) error {
	db := database.DB

	var stats struct {
		TotalProducts     int64 `json:"total_products"`
		TotalPortfolio    int64 `json:"total_portfolio"`
		TotalArticles     int64 `json:"total_articles"`
		TotalTestimonials int64 `json:"total_testimonials"`
		TotalOrders       int64 `json:"total_orders"`
		TotalChatSessions int64 `json:"total_chat_sessions"`
	}

	db.Model(&domain.Product{}).Count(&stats.TotalProducts)
	db.Model(&domain.Portfolio{}).Count(&stats.TotalPortfolio)
	db.Model(&domain.Article{}).Count(&stats.TotalArticles)
	db.Model(&domain.Testimonial{}).Count(&stats.TotalTestimonials)
	db.Model(&domain.OrderLog{}).Count(&stats.TotalOrders)
	db.Model(&domain.ChatSession{}).Count(&stats.TotalChatSessions)

	return c.JSON(stats)
}

// Publik — visitor ambil history chat berdasarkan session_id
func (h *ChatHandler) GetSessionHistory(c *fiber.Ctx) error {
	messages, err := h.service.GetSessionMessages(c.Params("session_id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Sesi tidak ditemukan"})
	}
	return c.JSON(messages)
}