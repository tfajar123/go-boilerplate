package route

import (
	"go-boilerplate/apps/internal/utils"
	"go-boilerplate/ent"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, client *ent.Client) {
	app.Get("/", func(c *fiber.Ctx) error {
		utils.Ok(c, "Service is running", nil)
		return nil
	})

	api := app.Group("/api/v1")

	registerAuthRoutes(api, client)

}
