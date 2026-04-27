package route

import (
	"go-boilerplate/apps/internal/config"
	"go-boilerplate/apps/internal/database"
	authController "go-boilerplate/apps/internal/features/auth/controllers"
	authService "go-boilerplate/apps/internal/features/auth/services"
	mailerService "go-boilerplate/apps/internal/features/mailer/services"
	storageService "go-boilerplate/apps/internal/features/storage/services"
	middlewares "go-boilerplate/apps/internal/middleware"
	"go-boilerplate/ent"
	"time"

	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(api fiber.Router, client *ent.Client, storage *database.Storage) {
	cfg := config.Load()
	storageSvc := storageService.NewStorageService(storage)
	mailSvc := mailerService.NewMailerService(cfg.Mail)
	authSvc := authService.NewAuthService(client, database.Redis, storageSvc, mailSvc)
	authCont := authController.NewAuthHandler(authSvc)
	redisClient := database.Redis

	api.Post("/auth/login", middlewares.RateLimiter(5, 1*time.Minute), authCont.Login)
	api.Post("/auth/register", authCont.Register)
	api.Post("/auth/forgot-password", middlewares.RateLimiter(3, 15*time.Minute), authCont.ForgotPassword)
	api.Post("/auth/reset-password", middlewares.RateLimiter(5, 15*time.Minute), authCont.ResetPassword)
	api.Post("/auth/refresh", authCont.Refresh)
	api.Post("/auth/logout", middlewares.AuthRequired(redisClient), authCont.Logout)
}
