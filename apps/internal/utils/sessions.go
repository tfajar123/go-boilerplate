package utils

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Session struct {
	UserID string
	Email  string
}

func CreateSession(
	ctx context.Context,
	rdb *redis.Client,
	userID string,
	email string,
	ttl time.Duration,
) (string, error) {

	sessionID := uuid.NewString()

	key := "session:" + sessionID

	err := rdb.HSet(ctx, key, map[string]any{
		"user_id": userID,
		"email":   email,
	}).Err()
	if err != nil {
		return "", err
	}

	rdb.Expire(ctx, key, ttl)

	return sessionID, nil
}

func GetSession(
	ctx context.Context,
	rdb *redis.Client,
	sessionID string,
) (*Session, error) {

	key := "session:" + sessionID

	data, err := rdb.HGetAll(ctx, key).Result()
	if err != nil || len(data) == 0 {
		return nil, redis.Nil
	}

	return &Session{
		UserID: data["user_id"],
		Email:  data["email"],
	}, nil
}

func DeleteSession(
	ctx context.Context,
	rdb *redis.Client,
	sessionID string,
) error {
	return rdb.Del(ctx, "session:"+sessionID).Err()
}
