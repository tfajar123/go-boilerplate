package authhandler

import (
	authService "go-boilerplate/apps/internal/features/auth/services"
	authValidation "go-boilerplate/apps/internal/features/auth/validation"
	"go-boilerplate/apps/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *authService.AuthService
}

func NewAuthHandler(authService *authService.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req authValidation.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid Body Request", err.Error())
	}

	userData, err := h.authService.Login(
		c.Context(),
		req,
	)
	if err != nil {

		if validationErrs := authValidation.FormatValidationError(err); len(validationErrs) > 0 {
			return utils.BadRequest(c, "Validation Failed", validationErrs)
		}

		return utils.Unauthorized(c, "Login Failed", err.Error())
	}

	return utils.Ok(c, "Login Success", userData)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req authValidation.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid Body Request", err.Error())
	}

	err := h.authService.Register(
		c.Context(),
		req,
	)
	if err != nil {

		if validationErrs := authValidation.FormatValidationError(err); len(validationErrs) > 0 {
			return utils.BadRequest(c, "Validation Failed", validationErrs)
		}

		return utils.BadRequest(c, "Registration Failed", err.Error())
	}

	return utils.Created(c, "Registration Success", nil)
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req authValidation.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid Body Request", err.Error())
	}

	result, err := h.authService.ForgotPassword(c.Context(), req)
	if err != nil {
		if validationErrs := authValidation.FormatValidationError(err); len(validationErrs) > 0 {
			return utils.BadRequest(c, "Validation Failed", validationErrs)
		}

		return utils.BadRequest(c, "Forgot Password Failed", err.Error())
	}

	return utils.Ok(c, "Forgot Password Success", result)
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req authValidation.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid Body Request", err.Error())
	}

	result, err := h.authService.ResetPassword(c.Context(), req)
	if err != nil {
		if validationErrs := authValidation.FormatValidationError(err); len(validationErrs) > 0 {
			return utils.BadRequest(c, "Validation Failed", validationErrs)
		}

		return utils.BadRequest(c, "Reset Password Failed", err.Error())
	}

	return utils.Ok(c, "Reset Password Success", result)
}

func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var req authValidation.VerifyRegisterOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid Body Request", err.Error())
	}

	result, err := h.authService.VerifyRegisterOTP(c.Context(), req)
	if err != nil {
		if validationErrs := authValidation.FormatValidationError(err); len(validationErrs) > 0 {
			return utils.BadRequest(c, "Validation Failed", validationErrs)
		}

		return utils.BadRequest(c, "Verify OTP Failed", err.Error())
	}

	return utils.Ok(c, "Verify OTP Success", result)
}

/* =========================
   REFRESH TOKEN
========================= */

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	type Req struct {
		RefreshToken string `json:"refresh_token"`
	}

	var req Req
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid Body Request", err.Error())
	}

	refreshToken, err := h.authService.RefreshToken(
		c.Context(),
		req.RefreshToken,
	)
	if err != nil {
		return utils.Unauthorized(c, "Invalid Token", err.Error())
	}

	return utils.Ok(c, "Token Refreshed", refreshToken)
}

/* =========================
   LOGOUT
========================= */

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Invalid session", nil)
	}
	sessionID, ok := c.Locals("session_id").(string)
	if !ok || sessionID == "" {
		return utils.Unauthorized(c, "Invalid session", nil)
	}

	if err := h.authService.Logout(
		c.Context(),
		userID,
		sessionID,
	); err != nil {
		return utils.Unauthorized(c, "Invalid Token", err.Error())
	}

	return utils.Ok(c, "Logout Success", nil)
}
