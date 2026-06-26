package handler

import (
	"BE/internal/service"
	"github.com/gofiber/fiber/v2"
)

type PortfolioHandler struct {
	service *service.PortfolioService
}

func NewPortfolioHandler(service *service.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{service: service}
}

func (h *PortfolioHandler) GetAll(c *fiber.Ctx) error {
	portfolios, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(portfolios)
}

func (h *PortfolioHandler) GetByID(c *fiber.Ctx) error {
	p, err := h.service.GetByID(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Karya tidak ditemukan"})
	}
	return c.JSON(p)
}

func (h *PortfolioHandler) Create(c *fiber.Ctx) error {
	var body map[string]any
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}
	p, err := h.service.Create(body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(201).JSON(p)
}

func (h *PortfolioHandler) Update(c *fiber.Ctx) error {
	var body map[string]any
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}
	p, err := h.service.Update(c.Params("id"), body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(p)
}

func (h *PortfolioHandler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.Params("id")); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Karya berhasil dihapus"})
}