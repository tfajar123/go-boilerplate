package authController

import (
	authService "go-boilerplate/apps/internal/features/auth/services"
	authValidation "go-boilerplate/apps/internal/features/auth/validation"
	"go-boilerplate/apps/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService *authService.AuthService
}

func NewAuthHandler(authService *authService.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (h *AuthController) Login(c *fiber.Ctx) error {
	var req authValidation.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "request body tidak valid", err.Error())
	}

	userData, accessToken, refreshToken, err := h.authService.Login(
		c.Context(),
		req,
	)
	if err != nil {

		if validationErrs := authValidation.FormatValidationError(err); len(validationErrs) > 0 {
			return utils.BadRequest(c, "validasi gagal", validationErrs)
		}

		return utils.Unauthorized(c, "login gagal", err.Error())
	}

	return utils.Ok(c, "login berhasil", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          userData,
	})
}

func (h *AuthController) Register(c *fiber.Ctx) error {
	var req authValidation.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "request body tidak valid", err.Error())
	}

	err := h.authService.Register(
		c.Context(),
		req,
	)
	if err != nil {

		if validationErrs := authValidation.FormatValidationError(err); len(validationErrs) > 0 {
			return utils.BadRequest(c, "validasi gagal", validationErrs)
		}

		return utils.BadRequest(c, "registrasi gagal", err.Error())
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
		return utils.BadRequest(c, "request body tidak valid", err.Error())
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(
		c.Context(),
		req.RefreshToken,
	)
	if err != nil {
		return utils.Unauthorized(c, "token tidak valid", err.Error())
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
		return utils.Unauthorized(c, "logout gagal", err.Error())
	}

	return utils.Ok(c, "logout berhasil", nil)
}
