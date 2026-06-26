package handler

import (
	"BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Request tidak valid"})
	}
	if body.Email == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Email dan password wajib diisi"})
	}

	token, user, err := h.authService.Login(body.Email, body.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}