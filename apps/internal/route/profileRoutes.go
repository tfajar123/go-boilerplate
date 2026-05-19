package route

import (
	"go-boilerplate/apps/internal/database"
	"go-boilerplate/apps/internal/features/profile/handlers"
	"go-boilerplate/apps/internal/features/profile/services"
	storageService "go-boilerplate/apps/internal/features/storage/services"
	middlewares "go-boilerplate/apps/internal/middleware"
	"go-boilerplate/ent"

	"github.com/gofiber/fiber/v2"
)

func registerProfileRoutes(api fiber.Router, client *ent.Client, storage *database.Storage) {
	storageSvc := storageService.NewStorageService(storage)
	profileSvc := services.NewProfileService(client, storageSvc)
	profileHndlr := handlers.NewProfileHandler(profileSvc)
	redisClient := database.Redis

	api.Get("/profiles", middlewares.AuthRequired(redisClient), profileHndlr.GetProfile)
	api.Put("/profiles/image", middlewares.AuthRequired(redisClient), profileHndlr.UpdateImageUrl)
}
