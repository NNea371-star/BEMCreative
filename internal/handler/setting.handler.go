package handler

import (
	"BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type SettingHandler struct {
	service *service.SettingService
}

func NewSettingHandler(service *service.SettingService) *SettingHandler {
	return &SettingHandler{service: service}
}

func (h *SettingHandler) GetPublicSettings(c *fiber.Ctx) error {
	settings, err := h.service.GetPublicSettings()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(settings)
}

func (h *SettingHandler) UpdateSettings(c *fiber.Ctx) error {
	var body struct {
		Settings []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"settings"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}

	items := make([]struct {
		Key   string
		Value string
	}, len(body.Settings))
	for i, s := range body.Settings {
		items[i].Key = s.Key
		items[i].Value = s.Value
	}

	if err := h.service.UpdateSettings(items); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Pengaturan berhasil disimpan"})
}

func (h *SettingHandler) GetWhatsappConfig(c *fiber.Ctx) error {
	config, err := h.service.GetWhatsappConfig()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(config)
}

func (h *SettingHandler) UpdateWhatsappConfig(c *fiber.Ctx) error {
	var body struct {
		PhoneNumber   string `json:"phone_number"`
		OrderTemplate string `json:"order_template"`
		HireTemplate  string `json:"hire_template"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}

	if err := h.service.UpdateWhatsappConfig(
		body.PhoneNumber,
		body.OrderTemplate,
		body.HireTemplate,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Konfigurasi WhatsApp berhasil disimpan"})
}

func (h *SettingHandler) GetPublicWhatsapp(c *fiber.Ctx) error {
	config, err := h.service.GetWhatsappConfig()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{
		"phone_number": config.PhoneNumber,
	})
}