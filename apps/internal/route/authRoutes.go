package route

import (
	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(api fiber.Router) {
	api.Post("/login", func(c *fiber.Ctx) error {
		return nil
	})
}