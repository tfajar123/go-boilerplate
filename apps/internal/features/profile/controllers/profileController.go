package controllers

import (
	"go-boilerplate/apps/internal/features/profile/services"
	"go-boilerplate/apps/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type ProfileController struct {
	profileService *services.ProfileService
}

func NewProfileController(profileService *services.ProfileService) *ProfileController {
	return &ProfileController{
		profileService: profileService,
	}
}

func (c *ProfileController) GetProfile(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("user_id").(string)

	if !ok || userId == "" {
		return utils.Unauthorized(ctx, "Unauthorized", ok)
	}

	profile, err := c.profileService.GetProfile(ctx.Context(), userId)
	if err != nil {
		return err
	}

	return utils.Ok(ctx, "Profile Fetched Successfully", profile)
}
