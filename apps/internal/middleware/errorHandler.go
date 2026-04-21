package middlewares

import (
	"go-boilerplate/apps/internal/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// GlobalErrorHandler is Fiber middleware for handling all errors globally
func GlobalErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err == nil {
			return nil
		}

		// Log original error
		utils.Logger.Error("Request error",
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.String("ip", c.IP()),
			zap.Error(err),
		)

		// Check if it's Fiber error
		fiberErr, isFiberErr := err.(*fiber.Error)
		if isFiberErr {
			return utils.Error(c, fiberErr.Code, fiberErr.Message, nil)
		}

		// Check if it's AppError
		statusCode := utils.GetStatusCode(err)
		message := utils.GetErrorMessage(err)

		// For production mode, don't expose internal error details
		// TODO: Check environment config
		var errorDetails any = nil
		if statusCode >= 500 {
			// Hide internal error details from client
			errorDetails = nil
		} else {
			// Send error details for client side errors
			errorDetails = err.Error()
		}

		return utils.Error(c, statusCode, message, errorDetails)
	}
}

// DatabaseErrorRecovery middleware handles database connection errors
func DatabaseErrorRecovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err == nil {
			return nil
		}

		// Check if error is database connection error
		if utils.IsDatabaseError(err) {
			appErr := utils.WrapDatabaseError(err)
			if appErr.StatusCode == http.StatusServiceUnavailable {
				// Database is down, add retry header
				c.Set("Retry-After", "60")
			}
		}

		return err
	}
}
