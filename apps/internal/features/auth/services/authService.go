package authService

import (
	"context"
	"errors"
	"time"

	"go-boilerplate/apps/internal/database"
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

func (s *AuthService) Login(ctx context.Context, email, password string) (*ent.User, string, string, error) {

	failKey := "auth:login_fail:" + email

	failCount, _ := s.redis.Get(ctx, failKey).Int()
	if failCount >= maxLoginAttempt {
		return nil, "", "", errors.New("terlalu banyak percobaan login, coba lagi nanti")
	}

	u, err := s.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(ctx)

	if err != nil {
		s.incrLoginFail(ctx, failKey)
		return nil, "", "", errors.New("email atau password salah")
	}

	if !utils.VerifyPassword(password, u.Password) {
		s.incrLoginFail(ctx, failKey)
		return nil, "", "", errors.New("email atau password salah")
	}

	// reset counter
	s.redis.Del(ctx, failKey)

	sid := utils.NewSessionID()

	err = s.SaveSession(ctx, u.ID.String(), sid, 7*24*time.Hour)
	if err != nil {
		return nil, "", "", err
	}

	accessToken, _ := utils.GenerateAccessToken(u.ID.String(), u.Email, sid)
	refreshToken, _ := utils.GenerateRefreshToken(u.ID.String(), u.Email, sid)

	return u, accessToken, refreshToken, nil
}

func (s *AuthService) Register(
	ctx context.Context,
	name string,
	email string,
	password string,
	role user.Role,
) error {

	exists, err := s.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Exist(ctx)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("email sudah terdaftar")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.client.User.
		Create().
		SetName(name).
		SetEmail(email).
		SetPassword(string(hashedPassword)).
		SetRole(role).
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
