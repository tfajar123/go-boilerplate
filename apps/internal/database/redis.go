package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"go-boilerplate/apps/internal/config"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func InitRedis(cfg config.RedisConfig) {
	opt := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:           cfg.DB,
		Username:     cfg.Username,
		Password:     cfg.Password,
		PoolSize:     20,
		MinIdleConns: 5,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		DialTimeout:  5 * time.Second,
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
