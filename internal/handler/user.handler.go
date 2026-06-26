package handler

import (
	"BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GET /api/admin/account — ambil profil user yang sedang login
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	user, err := h.userService.GetByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}

// PUT /api/admin/account — update profil + ganti password
func (h *UserHandler) UpdateAccount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var body struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}

	// Update profil
	user, err := h.userService.UpdateProfile(userID, body.Name, body.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	// Ganti password jika diisi
	if body.CurrentPassword != "" && body.NewPassword != "" {
		if err := h.userService.ChangePassword(userID, body.CurrentPassword, body.NewPassword); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": err.Error()})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Akun berhasil diperbarui",
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}