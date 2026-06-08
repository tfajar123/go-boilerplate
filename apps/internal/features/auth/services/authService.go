package authService

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"time"

	dto "go-boilerplate/apps/internal/features/auth/dto"
	authValidation "go-boilerplate/apps/internal/features/auth/validation"
	mailerService "go-boilerplate/apps/internal/features/mailer/services"
	storageService "go-boilerplate/apps/internal/features/storage/services"
	"go-boilerplate/apps/internal/utils"
	"go-boilerplate/models"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthService struct {
	db      *mongo.Database
	redis   *redis.Client
	storage *storageService.StorageService
	mailer  *mailerService.MailerService
}

func NewAuthService(
	db *mongo.Database,
	redis *redis.Client,
	storage *storageService.StorageService,
	mailer *mailerService.MailerService,
) *AuthService {
	return &AuthService{
		db:      db,
		redis:   redis,
		storage: storage,
		mailer:  mailer,
	}
}

const (
	maxLoginAttempt   = 5
	loginTTL          = 15 * time.Minute
	forgotPasswordTTL = 15 * time.Minute
	registerOTPTTL    = 15 * time.Minute
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

	var u models.User
	err := s.db.Collection(models.CollectionUsers).
		FindOne(ctx, bson.M{"email": req.Email}).
		Decode(&u)

	if err != nil {
		s.incrLoginFail(ctx, failKey)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("Email or Password is incorrect")
		}
		return nil, errors.New("Email or Password is incorrect")
	}

	if !utils.VerifyPassword(u.Password, req.Password) {
		s.incrLoginFail(ctx, failKey)
		return nil, errors.New("Email or Password is incorrect")
	}

	s.redis.Del(ctx, failKey)

	sid := utils.NewSessionID()

	err = s.SaveSession(ctx, u.ID.Hex(), sid, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	accessToken, _ := utils.GenerateAccessToken(u.ID.Hex(), u.Email, sid)
	refreshToken, _ := utils.GenerateRefreshToken(u.ID.Hex(), u.Email, sid)

	userResponse := dto.UserResponse{
		ID:        u.ID.Hex(),
		Name:      u.Profile.Name,
		Email:     u.Email,
		Role:      u.Role,
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

	count, err := s.db.Collection(models.CollectionUsers).
		CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("Email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// Generate 6 digit OTP
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	otpKey := "auth:register:otp:" + req.Email

	// Store registration data in Redis
	err = s.redis.HSet(ctx, otpKey, map[string]any{
		"email":    req.Email,
		"password": string(hashedPassword),
		"name":     req.Name,
		"role":     "user",
		"otp":      otp,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to store registration data: %w", err)
	}

	err = s.redis.Expire(ctx, otpKey, registerOTPTTL).Err()
	if err != nil {
		return err
	}

	// Development fallback: print OTP to terminal log
	log.Printf("\n\n📧 REGISTER OTP for %s (%s): %s (expires in %v)\n\n", req.Name, req.Email, otp, registerOTPTTL)

	// if s.mailer != nil {
	// 	if err := s.mailer.SendRegisterOTPEmail(req.Email, req.Name, otp, registerOTPTTL); err != nil {
	// 		// _ = s.redis.Del(ctx, otpKey).Err()
	// 		// return fmt.Errorf("failed to send OTP email: %w", err)
	// 		log.Printf("⚠️  Failed to send OTP email, falling back to terminal log: %v", err)
	// 	}
	// }

	return nil
}

func (s *AuthService) VerifyRegisterOTP(
	ctx context.Context,
	req authValidation.VerifyRegisterOTPRequest,
) (*dto.UserResponse, error) {

	if err := authValidation.ValidateAuth(req); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	otpKey := "auth:register:otp:" + req.Email

	storedOTP, err := s.redis.HGet(ctx, otpKey, "otp").Result()
	if err == redis.Nil {
		return nil, errors.New("OTP has expired or invalid")
	}
	if err != nil {
		return nil, err
	}

	if storedOTP != req.OTP {
		return nil, errors.New("Invalid OTP code")
	}

	// Get all registration data
	data, err := s.redis.HGetAll(ctx, otpKey).Result()
	if err != nil {
		return nil, err
	}

	// Create user document with embedded profile
	now := time.Now()
	user := models.User{
		Email:         data["email"],
		Password:      data["password"],
		Role:          data["role"],
		EmailVerified: true,
		Profile: models.Profile{
			Name: data["name"],
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := s.db.Collection(models.CollectionUsers).InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	insertedID := result.InsertedID.(bson.ObjectID)

	// Delete OTP
	_ = s.redis.Del(ctx, otpKey).Err()

	// Auto login after successful verification
	sid := utils.NewSessionID()
	err = s.SaveSession(ctx, insertedID.Hex(), sid, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        insertedID.Hex(),
		Name:      data["name"],
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *AuthService) incrLoginFail(
	ctx context.Context,
	key string,
) {
	count, err := s.redis.Incr(ctx, key).Result()
	if err != nil {
		return
	}

	if count == 1 {
		s.redis.Expire(ctx, key, loginTTL)
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

	var u models.User
	err := s.db.Collection(models.CollectionUsers).
		FindOne(ctx, bson.M{"email": req.Email}).
		Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("Email not found")
		}
		return nil, err
	}

	resetToken := bson.NewObjectID().Hex()
	expiresAt := time.Now().Add(forgotPasswordTTL)

	tokenKey := "auth:forgot_password:token:" + resetToken
	emailKey := "auth:forgot_password:email:" + req.Email

	existingToken, err := s.redis.Get(ctx, emailKey).Result()
	if err == nil && existingToken != "" {
		_ = s.redis.Del(ctx, "auth:forgot_password:token:"+existingToken).Err()
	} else if err != nil && err != redis.Nil {
		return nil, err
	}

	if err := s.redis.Set(ctx, tokenKey, u.ID.Hex(), forgotPasswordTTL).Err(); err != nil {
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
	if err := s.mailer.SendResetPasswordEmail(u.Email, u.Profile.Name, resetLink, expiresAt); err != nil {
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

	parsedID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in reset token: %w", err)
	}

	var u models.User
	err = s.db.Collection(models.CollectionUsers).
		FindOne(ctx, bson.M{"_id": parsedID}).
		Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Collection(models.CollectionUsers).
		UpdateByID(ctx, u.ID, bson.M{
			"$set": bson.M{
				"password":   hashedPassword,
				"updated_at": time.Now(),
			},
		})
	if err != nil {
		return nil, err
	}

	emailKey := "auth:forgot_password:email:" + u.Email
	sessionKey := "auth:session:" + u.ID.Hex()

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
