package handler

import (
	"BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// Public
func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	products, err := h.productService.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(products)
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	product, err := h.productService.GetByID(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Produk tidak ditemukan"})
	}
	return c.JSON(product)
}

func (h *ProductHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.productService.GetAllCategories()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(categories)
}

// Admin
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var body map[string]any
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}
	if body["name"] == nil || body["category_id"] == nil {
		return c.Status(400).JSON(fiber.Map{"message": "Nama dan kategori wajib diisi"})
	}

	product, err := h.productService.Create(body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(201).JSON(product)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	var body map[string]any
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}

	product, err := h.productService.Update(c.Params("id"), body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(product)
}

func (h *ProductHandler) UpdateStock(c *fiber.Ctx) error {
	var body struct {
		Stock int `json:"stock"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}

	if err := h.productService.UpdateStock(c.Params("id"), body.Stock); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Stok berhasil diperbarui"})
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	if err := h.productService.Delete(c.Params("id")); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Produk berhasil dihapus"})
}

func (h *ProductHandler) CreateCategory(c *fiber.Ctx) error {
	var body struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&body); err != nil || body.Name == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Nama kategori wajib diisi"})
	}

	cat, err := h.productService.CreateCategory(body.Name)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(201).JSON(cat)
}

func (h *ProductHandler) DeleteCategory(c *fiber.Ctx) error {
	if err := h.productService.DeleteCategory(c.Params("id")); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Kategori berhasil dihapus"})
}