package route

import (
	"go-boilerplate/apps/internal/database"
	authController "go-boilerplate/apps/internal/features/auth/controllers"
	authService "go-boilerplate/apps/internal/features/auth/services"
	middlewares "go-boilerplate/apps/internal/middleware"
	"go-boilerplate/ent"
	"time"

	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(api fiber.Router, client *ent.Client) {
	authSvc := authService.NewAuthService(client, database.Redis)
	authCont := authController.NewAuthHandler(authSvc)
	redisClient := database.Redis

	api.Post("/auth/login", middlewares.RateLimiter(5, 1*time.Minute), authCont.Login)
	api.Post("/auth/register", authCont.Register)
	api.Post("/auth/refresh", authCont.Refresh)
	api.Post("/auth/logout", middlewares.AuthRequired(redisClient), authCont.Logout)
}
