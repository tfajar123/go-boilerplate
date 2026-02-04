package main

import (
	"context"
	"go-boilerplate/apps/internal/config"
	"go-boilerplate/apps/internal/database"
	middlewares "go-boilerplate/apps/internal/middleware"
	"go-boilerplate/apps/internal/route"
	"go-boilerplate/apps/internal/utils"

	"github.com/gofiber/fiber/v2"
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

	minioClient := database.NewMinioClient(cfg.Minio)
	database.EnsureBucket(context.Background(), minioClient, cfg.Minio.Bucket)

	// init fiber
	app := fiber.New()
	app.Use(middlewares.RequestLogger())

	// inject dependencies
	route.Register(app, entClient)

	app.Listen(":" + cfg.Port)
}
