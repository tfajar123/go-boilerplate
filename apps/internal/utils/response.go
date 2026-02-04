package utils

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type SuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   any    `json:"error"`
}

func Success(c *fiber.Ctx, status int, message string, data any) error {
	resp := SuccessResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
	return c.Status(status).JSON(resp)

}

func Error(c *fiber.Ctx, status int, message string, err any) error {
	resp := ErrorResponse{
		Status:  status,
		Message: message,
		Error:   err,
	}
	return c.Status(status).JSON(resp)

}

// Shortcut
func Ok(c *fiber.Ctx, message string, data any) error {
	return Success(c, http.StatusOK, message, data)
}

func Created(c *fiber.Ctx, message string, data any) error {
	return Success(c, http.StatusCreated, message, data)
}

func NoContent(c *fiber.Ctx) error {
	return Success(c, http.StatusNoContent, "", nil)
}

func NotFound(c *fiber.Ctx, message string, err any) error {
	return Error(c, http.StatusNotFound, message, err)
}

func BadRequest(c *fiber.Ctx, message string, err any) error {
	return Error(c, http.StatusBadRequest, message, err)
}

func InternalError(c *fiber.Ctx, message string, err any) error {
	return Error(c, http.StatusInternalServerError, message, err)
}

func Unauthorized(c *fiber.Ctx, message string, err any) error {
	return Error(c, http.StatusUnauthorized, message, err)
}
