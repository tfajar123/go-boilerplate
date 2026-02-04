package authService

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-boilerplate/apps/internal/database"
	authValidation "go-boilerplate/apps/internal/features/auth/validation"
	"go-boilerplate/apps/internal/utils"
	"go-boilerplate/ent"
	"go-boilerplate/ent/user"

	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	client *ent.Client
	redis  *redis.Client
}

func NewAuthService(client *ent.Client, redis *redis.Client) *AuthService {
	return &AuthService{
		client: client,
		redis:  redis,
	}
}

const (
	maxLoginAttempt = 5
	loginTTL        = 15 * time.Minute
)

func (s *AuthService) Login(
	ctx context.Context,
	req authValidation.LoginRequest,
) (*ent.User, string, string, error) {

	if err := authValidation.ValidateAuth(req); err != nil {
		return nil, "", "", err
	}

	failKey := "auth:login_fail:" + req.Email

	failCount, _ := s.redis.Get(ctx, failKey).Int()
	if failCount >= maxLoginAttempt {
		return nil, "", "", errors.New("terlalu banyak percobaan login, coba lagi nanti")
	}

	u, err := s.client.User.
		Query().
		Where(user.EmailEQ(req.Email)).
		Only(ctx)

	if err != nil {
		s.incrLoginFail(ctx, failKey)
		return nil, "", "", errors.New("email atau password salah")
	}

	if !utils.VerifyPassword(u.Password, req.Password) {
		s.incrLoginFail(ctx, failKey)
		return nil, "", "", errors.New("email atau password salah")
	}

	s.redis.Del(ctx, failKey)

	sid := utils.NewSessionID()

	err = s.SaveSession(ctx, u.ID.String(), sid, 7*24*time.Hour)
	if err != nil {
		return nil, "", "", err
	}
	u.Password = ""

	accessToken, _ := utils.GenerateAccessToken(u.ID.String(), u.Email, sid)
	refreshToken, _ := utils.GenerateRefreshToken(u.ID.String(), u.Email, sid)

	return u, accessToken, refreshToken, nil
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
		return errors.New("email sudah terdaftar")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	_, err = s.client.User.
		Create().
		SetName(req.Name).
		SetEmail(req.Email).
		SetPassword(string(hashedPassword)).
		SetRole(req.Role).
		Save(ctx)

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
) (string, string, error) {

	claims, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	userID := claims["sub"].(string)
	oldSID := claims["sid"].(string)

	key := "auth:session:" + userID

	storedSID, err := s.redis.Get(ctx, key).Result()
	if err != nil || storedSID != oldSID {
		return "", "", errors.New("session tidak valid")
	}

	newSID := utils.NewSessionID()

	err = s.SaveSession(ctx, userID, newSID, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	accessToken, _ := utils.GenerateAccessToken(userID, claims["email"].(string), newSID)
	refreshTokenNew, _ := utils.GenerateRefreshToken(userID, claims["email"].(string), newSID)

	return accessToken, refreshTokenNew, nil
}

func (s *AuthService) Logout(
	ctx context.Context,
	userID string,
	sessionID string,
) error {

	key := "auth:session:" + userID

	storedSID, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return errors.New("sudah logout")
	}
	if err != nil {
		return err
	}

	if storedSID != sessionID {
		return errors.New("session tidak valid")
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
