package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-boilerplate/apps/internal/config"
	"go-boilerplate/apps/internal/database"
	middlewares "go-boilerplate/apps/internal/middleware"
	"go-boilerplate/apps/internal/route"
	"go-boilerplate/apps/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
)

func main() {
	utils.InitLogger()
	defer utils.SyncLogger()
	utils.Logger.Info("application starting")

	// load config
	cfg := config.Load()

	// init JWT secrets from config
	utils.InitJWT(cfg.JWTAccessSecret, cfg.JWTRefreshSecret)

	// init mongodb
	mongoClient, db := database.NewMongoClient(cfg.MongoURI, cfg.MongoDBName)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			utils.Logger.Error("failed to disconnect mongodb", zap.Error(err))
		}
	}()

	utils.Logger.Info("mongodb connected",
		zap.String("database", cfg.MongoDBName),
	)

	// bootstrap database (ensure indexes)
	database.Bootstrap(db)

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

	// init fiber with tuned config
	app := fiber.New(fiber.Config{
		ServerHeader:  "go-boilerplate",
		ReadTimeout:   10 * time.Second,
		WriteTimeout:  10 * time.Second,
		IdleTimeout:   120 * time.Second,
		BodyLimit:     10 * 1024 * 1024, // 10MB
		StrictRouting: false,
		CaseSensitive: false,
	})
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(helmet.New())
	app.Use(requestid.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CorsOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		AllowCredentials: cfg.CorsOrigins != "*",
	}))
	app.Use(middlewares.RequestLogger())

	// inject dependencies
	route.Register(app, db, storage, cfg)

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			utils.Logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	utils.Logger.Info("server started",
		zap.String("port", cfg.Port),
		zap.String("env", cfg.AppEnv),
	)

	<-quit
	utils.Logger.Info("shutting down server...")

	if err := app.Shutdown(); err != nil {
		utils.Logger.Error("server shutdown error", zap.Error(err))
	}

	utils.Logger.Info("server stopped")
}
