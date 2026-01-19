package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"go-boilerplate/apps/internal/config"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func InitRedis(cfg config.RedisConfig) {
	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:       cfg.DB,
		Username: cfg.Username,
		Password: cfg.Password,
	}

	if cfg.TLS {
		opt.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	Redis = redis.NewClient(opt)

	if err := Redis.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed connect redis: %v", err)
	}
}
