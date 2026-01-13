package handler

import (
	"go-boilerplate/apps/internal/service"

	"github.com/gofiber/fiber/v2"
)

func Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": service.HealthCheck(),
	})
}
