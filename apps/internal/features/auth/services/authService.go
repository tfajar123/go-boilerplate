package authService

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-boilerplate/apps/internal/database"
	dto "go-boilerplate/apps/internal/features/auth/dto"
	authValidation "go-boilerplate/apps/internal/features/auth/validation"
	storageService "go-boilerplate/apps/internal/features/storage/services"
	"go-boilerplate/apps/internal/utils"
	"go-boilerplate/ent"
	"go-boilerplate/ent/user"

	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	client  *ent.Client
	redis   *redis.Client
	storage *storageService.StorageService
}

func NewAuthService(
	client *ent.Client,
	redis *redis.Client,
	storage *storageService.StorageService,
) *AuthService {
	return &AuthService{
		client:  client,
		redis:   redis,
		storage: storage,
	}
}

const (
	maxLoginAttempt = 5
	loginTTL        = 15 * time.Minute
)

func (s *AuthService) Login(
	ctx context.Context,
	req authValidation.LoginRequest,
) (*dto.LoginResponse, error) {

	if err := authValidation.ValidateAuth(req); err != nil {
		return nil, err
	}

	failKey := "auth:login_fail:" + req.Email

	failCount, _ := s.redis.Get(ctx, failKey).Int()
	if failCount >= maxLoginAttempt {
		return nil, errors.New("Too many requests, please try again later")
	}

	u, err := s.client.User.
		Query().
		Where(user.EmailEQ(req.Email)).
		Only(ctx)

	if err != nil {
		s.incrLoginFail(ctx, failKey)
		return nil, errors.New("Email or Password is incorrect")
	}

	if !utils.VerifyPassword(u.Password, req.Password) {
		s.incrLoginFail(ctx, failKey)
		return nil, errors.New("Email or Password is incorrect")
	}

	s.redis.Del(ctx, failKey)

	sid := utils.NewSessionID()

	err = s.SaveSession(ctx, u.ID.String(), sid, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}
	u.Password = ""

	accessToken, _ := utils.GenerateAccessToken(u.ID.String(), u.Email, sid)
	refreshToken, _ := utils.GenerateRefreshToken(u.ID.String(), u.Email, sid)

	userResponse := dto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role.String(),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	loginResponse := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userResponse,
	}

	return &loginResponse, nil
}

func (s *AuthService) Register(
	ctx context.Context,
	req authValidation.RegisterRequest,
) error {

	if err := authValidation.ValidateAuth(req); err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	exists, err := s.client.User.
		Query().
		Where(user.EmailEQ(req.Email)).
		Exist(ctx)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("Email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// Skip upload jika image kosong
	var profileImage string
	if req.Image != "" {

		image, err := s.storage.UploadBase64(ctx, req.Image, "profiles")
		if err != nil {
			return fmt.Errorf("failed to upload profile image: %w", err)
		}
		profileImage = image
	}

	userCreate := s.client.User.
		Create().
		SetName(req.Name).
		SetEmail(req.Email).
		SetPassword(string(hashedPassword)).
		SetRole(req.Role)

	// Ent otomatis generate method camelCase sesuai nama field: profileImage -> SetProfileImage
	if profileImage != "" {
		userCreate.SetProfileImage(profileImage)
	}

	_, err = userCreate.Save(ctx)

	return err
}

func (s *AuthService) incrLoginFail(
	ctx context.Context,
	key string,
) {
	count, err := database.Redis.Incr(ctx, key).Result()
	if err != nil {
		return
	}

	if count == 1 {
		database.Redis.Expire(ctx, key, loginTTL)
	}
}

func (s *AuthService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*dto.RefreshTokenResponse, error) {

	claims, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userID := claims["sub"].(string)
	oldSID := claims["sid"].(string)

	key := "auth:session:" + userID

	storedSID, err := s.redis.Get(ctx, key).Result()
	if err != nil || storedSID != oldSID {
		return nil, errors.New("Invalid Session ID")
	}

	newSID := utils.NewSessionID()

	err = s.SaveSession(ctx, userID, newSID, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	accessToken, _ := utils.GenerateAccessToken(userID, claims["email"].(string), newSID)
	refreshTokenNew, _ := utils.GenerateRefreshToken(userID, claims["email"].(string), newSID)

	response := dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenNew,
	}

	return &response, nil
}

func (s *AuthService) Logout(
	ctx context.Context,
	userID string,
	sessionID string,
) error {

	key := "auth:session:" + userID

	storedSID, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return errors.New("Logged out")
	}
	if err != nil {
		return err
	}

	if storedSID != sessionID {
		return errors.New("Invalid Session ID")
	}

	return s.redis.Del(ctx, key).Err()
}

func (s *AuthService) SaveSession(
	ctx context.Context,
	userID string,
	sessionID string,
	ttl time.Duration,
) error {
	return s.redis.Set(
		ctx,
		"auth:session:"+userID,
		sessionID,
		ttl,
	).Err()
}
