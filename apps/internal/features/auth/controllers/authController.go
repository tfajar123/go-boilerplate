package authController

import (
	authService "go-boilerplate/apps/internal/features/auth/services"
	"go-boilerplate/apps/internal/utils"
	"go-boilerplate/ent/user"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService *authService.AuthService
}

func NewAuthHandler(authService *authService.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *AuthController) Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		utils.BadRequest(c, "request body tidak valid")
		return nil
	}

	if req.Email == "" || req.Password == "" {
		utils.BadRequest(c, "email dan password wajib diisi")
		return nil
	}

	user, err := h.authService.Login(
		c.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return nil
	}

	token, err := utils.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		utils.InternalError(c, "gagal membuat token")
		return nil
	}

	return utils.Ok(c, "login berhasil", map[string]any{"user": map[string]any{
		"id":    user.ID.String(),
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}, "token": token})
}

func (h *AuthController) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "request body tidak valid")
	}

	if err := h.authService.Register(
		c.Context(),
		req.Name,
		req.Email,
		req.Password,
		user.Role(req.Role),
	); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "registrasi berhasil", nil)
}
