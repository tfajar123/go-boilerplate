package controllers

import (
	"go-boilerplate/apps/internal/features/profile/services"
	"go-boilerplate/apps/internal/features/profile/validation"
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

func (c *ProfileController) UpdateImageUrl(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("user_id").(string)
	if !ok || userId == "" {
		return utils.Unauthorized(ctx, "Unauthorized", nil)
	}

	var req validation.UpdateProfileImageRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid Body Request", err.Error())
	}

	result, err := c.profileService.UpdateImageUrl(ctx.Context(), userId, req)
	if err != nil {
		if validationErrs := validation.FormatValidationError(err); len(validationErrs) > 0 {
			return utils.BadRequest(ctx, "Validation Failed", validationErrs)
		}

		return err
	}

	return utils.Ok(ctx, "Profile Image Updated Successfully", result)
}
