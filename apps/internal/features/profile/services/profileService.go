package services

import (
	"context"
	"errors"
	"fmt"
	"go-boilerplate/apps/internal/config"
	"go-boilerplate/apps/internal/features/profile/dto"
	profileValidation "go-boilerplate/apps/internal/features/profile/validation"
	storageService "go-boilerplate/apps/internal/features/storage/services"
	"go-boilerplate/apps/internal/utils"
	"go-boilerplate/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type ProfileService struct {
	db      *mongo.Database
	storage *storageService.StorageService
}

func NewProfileService(db *mongo.Database, storage *storageService.StorageService) *ProfileService {
	return &ProfileService{
		db:      db,
		storage: storage,
	}
}

func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*dto.ProfileResponse, error) {
	parsedID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	var u models.User
	err = s.db.Collection(models.CollectionUsers).
		FindOne(ctx, bson.M{"_id": parsedID}).
		Decode(&u)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	cfg := config.Load()
	profileImg := ""
	if u.Profile.ImageUrl != "" {
		profileImg = utils.BuildStorageURL(cfg, u.Profile.ImageUrl)
	}

	return &dto.ProfileResponse{
		ID:        u.ID.Hex(),
		Name:      u.Profile.Name,
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

	parsedID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	// Get current user to check for existing image
	var u models.User
	err = s.db.Collection(models.CollectionUsers).
		FindOne(ctx, bson.M{"_id": parsedID}).
		Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	imagePath, err := s.storage.UploadBase64(ctx, req.Image, "profiles")
	if err != nil {
		return nil, fmt.Errorf("failed to upload profile image: %w", err)
	}

	previousImage := u.Profile.ImageUrl

	// Update profile.image_url using $set on embedded field
	_, err = s.db.Collection(models.CollectionUsers).
		UpdateByID(ctx, parsedID, bson.M{
			"$set": bson.M{
				"profile.image_url": imagePath,
				"updated_at":        time.Now(),
			},
		})
	if err != nil {
		return nil, err
	}

	if previousImage != "" && previousImage != imagePath {
		if deleteErr := s.storage.Delete(ctx, previousImage); deleteErr != nil {
			utils.Logger.Warn("failed to delete old profile image",
				zap.String("user_id", parsedID.Hex()),
				zap.String("old_image", previousImage),
				zap.Error(deleteErr),
			)
		}
	}

	cfg := config.Load()

	return &dto.UpdateProfileImageResponse{
		ImageUrl: utils.BuildStorageURL(cfg, imagePath),
	}, nil
}
