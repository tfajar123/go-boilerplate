package route

import (
	"go-boilerplate/apps/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	app.Get("/health", handler.Health)
}
