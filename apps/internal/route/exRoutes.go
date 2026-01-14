package route

import (
	middlewares "go-boilerplate/apps/internal/middleware"
	"go-boilerplate/ent"

	"github.com/gofiber/fiber/v2"
)

func registerExRoutes(api fiber.Router, client *ent.Client) {
	protected := api.Group("/user", middlewares.AuthRequired())

	protected.Get("/profile", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user_id": c.Locals("user_id"),
			"email":   c.Locals("email"),
		})
	})
}