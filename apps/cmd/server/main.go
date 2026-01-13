package main

import (
	"go-boilerplate/apps/internal/config"
	"go-boilerplate/apps/internal/database"
	"go-boilerplate/apps/internal/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()

	// bootstrap database
	database.EnsureDatabaseExists(cfg.DBUrl)
	database.RunMigration(cfg.DBUrl)

	// init ent
	entClient := database.NewEntClient(cfg.DBUrl)
	defer entClient.Close()

	// init fiber
	app := fiber.New()

	// inject dependencies
	route.Register(app, entClient)

	app.Listen(":" + cfg.Port)
}

