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

/* =========================
   LOGIN
========================= */

func (h *AuthController) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "request body tidak valid")
	}

	userData, accessToken, refreshToken, err := h.authService.Login(
		c.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		return utils.Unauthorized(c, err.Error())
	}

	return utils.Ok(c, "login berhasil", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          userData,
	})
}

/* =========================
   REGISTER
========================= */

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

/* =========================
   REFRESH TOKEN
========================= */

func (h *AuthController) Refresh(c *fiber.Ctx) error {
	type Req struct {
		RefreshToken string `json:"refresh_token"`
	}

	var req Req
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "request body tidak valid")
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(
		c.Context(),
		req.RefreshToken,
	)
	if err != nil {
		return utils.Unauthorized(c, err.Error())
	}

	return utils.Ok(c, "token diperbarui", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

/* =========================
   LOGOUT
========================= */

func (h *AuthController) Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	sessionID := c.Locals("session_id").(string)

	if err := h.authService.Logout(
		c.Context(),
		userID,
		sessionID,
	); err != nil {
		return utils.Unauthorized(c, err.Error())
	}

	return utils.Ok(c, "logout berhasil", nil)
}
