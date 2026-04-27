package authService

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"go-boilerplate/apps/internal/database"
	dto "go-boilerplate/apps/internal/features/auth/dto"
	authValidation "go-boilerplate/apps/internal/features/auth/validation"
	"go-boilerplate/apps/internal/features/mailer/services"
	storageService "go-boilerplate/apps/internal/features/storage/services"
	"go-boilerplate/apps/internal/utils"
	"go-boilerplate/ent"
	"go-boilerplate/ent/user"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	client  *ent.Client
	redis   *redis.Client
	storage *storageService.StorageService
	mailer  *mailerService.MailerService
}

func NewAuthService(
	client *ent.Client,
	redis *redis.Client,
	storage *storageService.StorageService,
	mailer *mailerService.MailerService,
) *AuthService {
	return &AuthService{
		client:  client,
		redis:   redis,
		storage: storage,
		mailer:  mailer,
	}
}

const (
	maxLoginAttempt   = 5
	loginTTL          = 15 * time.Minute
	forgotPasswordTTL = 15 * time.Minute
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

func (s *AuthService) ForgotPassword(
	ctx context.Context,
	req authValidation.ForgotPasswordRequest,
) (*dto.ForgotPasswordResponse, error) {
	if err := authValidation.ValidateAuth(req); err != nil {
		return nil, err
	}

	u, err := s.client.User.
		Query().
		Where(user.EmailEQ(req.Email)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("Email not found")
		}
		return nil, err
	}

	resetToken := uuid.NewString()
	expiresAt := time.Now().Add(forgotPasswordTTL)

	tokenKey := "auth:forgot_password:token:" + resetToken
	emailKey := "auth:forgot_password:email:" + req.Email

	existingToken, err := s.redis.Get(ctx, emailKey).Result()
	if err == nil && existingToken != "" {
		_ = s.redis.Del(ctx, "auth:forgot_password:token:"+existingToken).Err()
	} else if err != nil && err != redis.Nil {
		return nil, err
	}

	if err := s.redis.Set(ctx, tokenKey, u.ID.String(), forgotPasswordTTL).Err(); err != nil {
		return nil, err
	}

	if err := s.redis.Set(ctx, emailKey, resetToken, forgotPasswordTTL).Err(); err != nil {
		return nil, err
	}

	if s.mailer == nil {
		_ = s.redis.Del(ctx, tokenKey, emailKey).Err()
		return nil, errors.New("mailer service is not configured")
	}

	resetLink := fmt.Sprintf("%s?token=%s", s.mailerResetPasswordURL(), url.QueryEscape(resetToken))
	if err := s.mailer.SendResetPasswordEmail(u.Email, u.Name, resetLink, expiresAt); err != nil {
		_ = s.redis.Del(ctx, tokenKey, emailKey).Err()
		return nil, fmt.Errorf("failed to send reset password email: %w", err)
	}

	return &dto.ForgotPasswordResponse{
		EmailSent: true,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) ResetPassword(
	ctx context.Context,
	req authValidation.ResetPasswordRequest,
) (*dto.ResetPasswordResponse, error) {
	if err := authValidation.ValidateAuth(req); err != nil {
		return nil, err
	}

	tokenKey := "auth:forgot_password:token:" + req.Token

	userID, err := s.redis.Get(ctx, tokenKey).Result()
	if err == redis.Nil {
		return nil, errors.New("Reset token is invalid or expired")
	}
	if err != nil {
		return nil, err
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in reset token: %w", err)
	}

	u, err := s.client.User.
		Query().
		Where(user.IDEQ(parsedID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, err
	}

	if err := s.client.User.
		UpdateOneID(u.ID).
		SetPassword(hashedPassword).
		Exec(ctx); err != nil {
		return nil, err
	}

	emailKey := "auth:forgot_password:email:" + u.Email
	sessionKey := "auth:session:" + u.ID.String()

	if err := s.redis.Del(ctx, tokenKey, emailKey, sessionKey).Err(); err != nil {
		return nil, err
	}

	return &dto.ResetPasswordResponse{
		ResetAt: time.Now(),
	}, nil
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

func (s *AuthService) mailerResetPasswordURL() string {
	if s.mailer == nil {
		return ""
	}

	return s.mailer.ResetPasswordURL()
}
