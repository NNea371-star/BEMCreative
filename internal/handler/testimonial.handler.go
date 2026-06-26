package handler

import (
	"BE/internal/service"
	"github.com/gofiber/fiber/v2"
)

type TestimonialHandler struct {
	service *service.TestimonialService
}

func NewTestimonialHandler(service *service.TestimonialService) *TestimonialHandler {
	return &TestimonialHandler{service: service}
}

func (h *TestimonialHandler) GetVisible(c *fiber.Ctx) error {
	testimonials, err := h.service.GetVisible()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(testimonials)
}

func (h *TestimonialHandler) SubmitPublic(c *fiber.Ctx) error {
	var body map[string]any
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}
	body["is_visible"] = false // default pending moderasi
	t, err := h.service.Create(body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "Testimoni berhasil dikirim, menunggu persetujuan admin.",
		"data":    t,
	})
}

func (h *TestimonialHandler) GetAll(c *fiber.Ctx) error {
	testimonials, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(testimonials)
}

func (h *TestimonialHandler) Update(c *fiber.Ctx) error {
	var body map[string]any
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}
	t, err := h.service.Update(c.Params("id"), body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(t)
}

func (h *TestimonialHandler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.Params("id")); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Testimoni berhasil dihapus"})
}

func (h *TestimonialHandler) Create(c *fiber.Ctx) error {
	var body map[string]any
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}
	t, err := h.service.Create(body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(201).JSON(t)
}