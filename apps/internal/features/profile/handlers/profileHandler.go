package handlers

import (
	"go-boilerplate/apps/internal/features/profile/services"
	"go-boilerplate/apps/internal/features/profile/validation"
	"go-boilerplate/apps/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

func (h *ProfileHandler) GetProfile(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("user_id").(string)

	if !ok || userId == "" {
		return utils.Unauthorized(ctx, "Unauthorized", ok)
	}

	profile, err := h.profileService.GetProfile(ctx.Context(), userId)
	if err != nil {
		return err
	}

	return utils.Ok(ctx, "Profile Fetched Successfully", profile)
}

func (h *ProfileHandler) UpdateImageUrl(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("user_id").(string)
	if !ok || userId == "" {
		return utils.Unauthorized(ctx, "Unauthorized", nil)
	}

	var req validation.UpdateProfileImageRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid Body Request", err.Error())
	}

	result, err := h.profileService.UpdateImageUrl(ctx.Context(), userId, req)
	if err != nil {
		if validationErrs := validation.FormatValidationError(err); len(validationErrs) > 0 {
			return utils.BadRequest(ctx, "Validation Failed", validationErrs)
		}

		return err
	}

	return utils.Ok(ctx, "Profile Image Updated Successfully", result)
}
