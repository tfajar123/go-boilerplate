package main

import (
	"go-boilerplate/apps/internal/config"
	"go-boilerplate/apps/internal/database"
	"go-boilerplate/apps/internal/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()

	// Run migration
	database.RunMigration(cfg.DBUrl)

	// Init DB
	_ = database.NewEntClient(cfg.DBUrl)

	app := fiber.New()

	route.Register(app)

	app.Listen(":" + cfg.Port)
}
