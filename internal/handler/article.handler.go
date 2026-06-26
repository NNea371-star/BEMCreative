package handler

import (
	"BE/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ArticleHandler struct {
	service *service.ArticleService
}

func NewArticleHandler(service *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

func (h *ArticleHandler) GetPublished(c *fiber.Ctx) error {
	articles, err := h.service.GetPublished()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(articles)
}

func (h *ArticleHandler) GetBySlug(c *fiber.Ctx) error {
	a, err := h.service.GetBySlug(c.Params("slug"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Artikel tidak ditemukan"})
	}
	return c.JSON(a)
}

func (h *ArticleHandler) GetCategories(c *fiber.Ctx) error {
	cats, err := h.service.GetAllCategories()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(cats)
}

// Admin
func (h *ArticleHandler) GetAll(c *fiber.Ctx) error {
	status := c.Query("status")
	articles, err := h.service.GetAll(status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(articles)
}

func (h *ArticleHandler) Create(c *fiber.Ctx) error {
	var body map[string]any
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}
	userID := c.Locals("user_id").(string)
	a, err := h.service.Create(body, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(201).JSON(a)
}

func (h *ArticleHandler) Update(c *fiber.Ctx) error {
	var body map[string]any
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}
	a, err := h.service.Update(c.Params("id"), body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(a)
}

func (h *ArticleHandler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.Params("id")); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Artikel berhasil dihapus"})
}

func (h *ArticleHandler) CreateCategory(c *fiber.Ctx) error {
	var body struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&body); err != nil || body.Name == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Nama kategori wajib diisi"})
	}
	cat, err := h.service.CreateCategory(body.Name)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(201).JSON(cat)
}

func (h *ArticleHandler) DeleteCategory(c *fiber.Ctx) error {
	if err := h.service.DeleteCategory(c.Params("id")); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Kategori berhasil dihapus"})
}