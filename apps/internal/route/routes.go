package route

import (
	"go-boilerplate/apps/internal/database"
	middlewares "go-boilerplate/apps/internal/middleware"
	"go-boilerplate/apps/internal/utils"
	"go-boilerplate/ent"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, client *ent.Client, storage *database.Storage) {

	app.Get("/", func(c *fiber.Ctx) error {
		utils.Ok(c, "Service is running", nil)
		return nil
	})

	api := app.Group("/api/v1")
	api.Use(middlewares.RateLimiter(100, time.Minute))

	registerAuthRoutes(api, client)
	registerExRoutes(api, client)

}
