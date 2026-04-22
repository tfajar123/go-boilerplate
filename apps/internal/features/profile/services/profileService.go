package services

import (
	"context"
	"fmt"
	"go-boilerplate/apps/internal/config"
	"go-boilerplate/apps/internal/features/profile/dto"
	storageService "go-boilerplate/apps/internal/features/storage/services"
	"go-boilerplate/apps/internal/utils"
	"go-boilerplate/ent"
	"go-boilerplate/ent/user"

	"github.com/google/uuid"
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
