package handler

import (
	"BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type StatsHandler interface {
	GetSiteStats(c *fiber.Ctx) error
}

type statsHandler struct {
	statsService service.StatsService
}

func NewStatsHandler(statsService service.StatsService) StatsHandler {
	return &statsHandler{
		statsService: statsService,
	}
}

func (h *statsHandler) GetSiteStats(c *fiber.Ctx) error {
	stats, err := h.statsService.GetSiteStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil statistik",
		})
	}
	return c.JSON(stats)
}