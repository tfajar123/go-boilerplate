package route

import (
	"go-boilerplate/apps/internal/config"
	"go-boilerplate/apps/internal/database"
	authHandler "go-boilerplate/apps/internal/features/auth/handlers"
	authService "go-boilerplate/apps/internal/features/auth/services"
	mailerService "go-boilerplate/apps/internal/features/mailer/services"
	storageService "go-boilerplate/apps/internal/features/storage/services"
	middlewares "go-boilerplate/apps/internal/middleware"
	"go-boilerplate/ent"
	"time"

	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(api fiber.Router, client *ent.Client, storage *database.Storage, cfg *config.Config) {
	storageSvc := storageService.NewStorageService(storage)
	mailSvc := mailerService.NewMailerService(cfg.Mail)
	authSvc := authService.NewAuthService(client, database.Redis, storageSvc, mailSvc)
	authHndlr := authHandler.NewAuthHandler(authSvc)
	redisClient := database.Redis

	api.Post("/auth/login", middlewares.RateLimiter(5, 1*time.Minute), authHndlr.Login)
	api.Post("/auth/register", authHndlr.Register)
	api.Post("/auth/forgot-password", middlewares.RateLimiter(3, 15*time.Minute), authHndlr.ForgotPassword)
	api.Post("/auth/reset-password", middlewares.RateLimiter(5, 15*time.Minute), authHndlr.ResetPassword)
	api.Post("/auth/verify-otp", authHndlr.VerifyOTP)
	api.Post("/auth/refresh", authHndlr.Refresh)
	api.Post("/auth/logout", middlewares.AuthRequired(redisClient), authHndlr.Logout)
}
