package handler

import (
	"BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	uploadService *service.UploadService
}

func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

// POST /api/admin/upload/:folder
// folder: products | portfolio | articles | settings
func (h *UploadHandler) Upload(c *fiber.Ctx) error {
	folder := c.Params("folder")

	// Validasi folder
	allowed := map[string]bool{
		"products":  true,
		"portfolio": true,
		"articles":  true,
		"settings":  true,
		"logos":     true,
	}
	if !allowed[folder] {
		return c.Status(400).JSON(fiber.Map{"message": "Folder tidak valid"})
	}

	// Ambil file dari form
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "File gambar wajib disertakan"})
	}

	// Buka file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal membaca file"})
	}
	defer file.Close()

	// Upload ke Cloudinary
	url, err := h.uploadService.UploadImage(file, fileHeader, folder)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"url":     url,
		"message": "Gambar berhasil diupload",
	})
}

// DELETE /api/admin/upload
func (h *UploadHandler) Delete(c *fiber.Ctx) error {
	var body struct {
		PublicID string `json:"public_id"`
	}
	if err := c.BodyParser(&body); err != nil || body.PublicID == "" {
		return c.Status(400).JSON(fiber.Map{"message": "public_id wajib diisi"})
	}

	if err := h.uploadService.DeleteImage(body.PublicID); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menghapus gambar"})
	}

	return c.JSON(fiber.Map{"message": "Gambar berhasil dihapus"})
}