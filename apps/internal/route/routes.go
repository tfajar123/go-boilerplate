package route

import (
	"go-boilerplate/apps/internal/config"
	"go-boilerplate/apps/internal/database"
	middlewares "go-boilerplate/apps/internal/middleware"
	"go-boilerplate/apps/internal/utils"
	"go-boilerplate/ent"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, client *ent.Client, storage *database.Storage, cfg *config.Config) {

	app.Get("/", func(c *fiber.Ctx) error {
		return utils.Ok(c, "Service is running", nil)
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		checks := fiber.Map{}

		// check database
		if err := database.DB.PingContext(c.Context()); err != nil {
			checks["db"] = "error: " + err.Error()
		} else {
			checks["db"] = "ok"
		}

		// check redis
		if err := database.Redis.Ping(c.Context()).Err(); err != nil {
			checks["redis"] = "error: " + err.Error()
		} else {
			checks["redis"] = "ok"
		}

		return utils.Ok(c, "healthy", checks)
	})

	api := app.Group("/api/v1")
	api.Use(middlewares.RateLimiter(100, time.Minute))

	registerAuthRoutes(api, client, storage, cfg)
	registerProfileRoutes(api, client, storage)

}
