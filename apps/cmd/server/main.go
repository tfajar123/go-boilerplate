package main

import (
	"go-boilerplate/apps/internal/config"
	"go-boilerplate/apps/internal/database"
	middlewares "go-boilerplate/apps/internal/middleware"
	"go-boilerplate/apps/internal/route"
	"go-boilerplate/apps/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func main() {
	utils.InitLogger()
	defer utils.SyncLogger()
	utils.Logger.Info("application starting")

	// load config
	cfg := config.Load()

	// bootstrap database
	database.EnsureDatabaseExists(cfg.DBUrl)

	// init ent
	entClient := database.NewEntClient(cfg.DBUrl)
	defer entClient.Close()

	database.InitRedis(cfg.Redis)

	utils.Logger.Info("redis connected",
		zap.String("host", cfg.Redis.Host),
		zap.Int("db", cfg.Redis.DB),
	)

	storage := database.NewStorage(cfg.Storage)
	utils.Logger.Info("storage connected",
		zap.String("endpoint", cfg.Storage.Endpoint),
		zap.String("bucket", cfg.Storage.Bucket),
	)

	// init fiber
	app := fiber.New()
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // CHANGE THIS IN PRODUCTION
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false, // IMPORTANT!! CHANGE THIS TO TRUE IN PRODUCTION, BUT MAKE SURE TO SET ALLOWORIGINS TO YOUR FRONTEND DOMAIN
	}))
	app.Use(middlewares.RequestLogger())

	// inject dependencies
	route.Register(app, entClient, storage)

	if err := app.Listen(":" + cfg.Port); err != nil {
		utils.Logger.Fatal("failed to start server", zap.Error(err))
	}
}
