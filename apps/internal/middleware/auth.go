package middlewares

import (
	"strings"

	"go-boilerplate/apps/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "authorization header tidak ada",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "format token tidak valid",
			})
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "token tidak valid",
			})
		}

		// simpan ke context
		c.Locals("userId", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}
