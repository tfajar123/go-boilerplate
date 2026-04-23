package services

import (
	"context"
	"fmt"
	"go-boilerplate/apps/internal/config"
	"go-boilerplate/apps/internal/features/profile/dto"
	profileValidation "go-boilerplate/apps/internal/features/profile/validation"
	storageService "go-boilerplate/apps/internal/features/storage/services"
	"go-boilerplate/apps/internal/utils"
	"go-boilerplate/ent"
	"go-boilerplate/ent/user"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProfileService struct {
	client  *ent.Client
	storage *storageService.StorageService
}

func NewProfileService(client *ent.Client, storage *storageService.StorageService) *ProfileService {
	return &ProfileService{
		client:  client,
		storage: storage,
	}
}

func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*dto.ProfileResponse, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid format: %w", err)
	}

	u, err := s.client.User.
		Query().
		Where(user.IDEQ(parsedID)).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	cfg := config.Load()
	profileImg := utils.BuildStorageURL(cfg, u.ProfileImage)

	return &dto.ProfileResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		ImageUrl:  profileImg,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (s *ProfileService) UpdateImageUrl(
	ctx context.Context,
	userID string,
	req profileValidation.UpdateProfileImageRequest,
) (*dto.UpdateProfileImageResponse, error) {
	if err := profileValidation.ValidateProfile(req); err != nil {
		return nil, err
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid format: %w", err)
	}

	u, err := s.client.User.
		Query().
		Where(user.IDEQ(parsedID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	imagePath, err := s.storage.UploadBase64(ctx, req.Image, "profiles")
	if err != nil {
		return nil, fmt.Errorf("failed to upload profile image: %w", err)
	}

	previousImage := u.ProfileImage

	updatedUser, err := s.client.User.
		UpdateOneID(parsedID).
		SetProfileImage(imagePath).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if previousImage != "" && previousImage != imagePath {
		if deleteErr := s.storage.Delete(ctx, previousImage); deleteErr != nil {
			utils.Logger.Warn("failed to delete old profile image",
				zap.String("user_id", parsedID.String()),
				zap.String("old_image", previousImage),
				zap.Error(deleteErr),
			)
		}
	}

	cfg := config.Load()

	return &dto.UpdateProfileImageResponse{
		ImageUrl: utils.BuildStorageURL(cfg, updatedUser.ProfileImage),
	}, nil
}
