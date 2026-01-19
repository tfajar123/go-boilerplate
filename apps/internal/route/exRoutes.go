package route

import (
	"go-boilerplate/apps/internal/database"
	middlewares "go-boilerplate/apps/internal/middleware"
	"go-boilerplate/ent"

	"github.com/gofiber/fiber/v2"
)

func registerExRoutes(api fiber.Router, client *ent.Client) {
	redisClient := database.Redis
	protected := api.Group("/user", middlewares.AuthRequired(redisClient))

	protected.Get("/profile", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user_id":    c.Locals("user_id"),
			"email":      c.Locals("email"),
			"session_id": c.Locals("session_id"),
		})
	})
}
