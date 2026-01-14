package route

import (
	authController "go-boilerplate/apps/internal/features/auth/controllers"
	authService "go-boilerplate/apps/internal/features/auth/services"
	"go-boilerplate/ent"

	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(api fiber.Router, client *ent.Client) {
	authSvc := authService.NewAuthService(client) 
	authCont := authController.NewAuthHandler(authSvc)
	
	api.Post("/auth/login", authCont.Login)
	api.Post("/auth/register", authCont.Register)
}