package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type RedisConfig struct {
	Host     string
	Port     int
	DB       int
	Username string
	Password string
	TLS      bool
}

type Config struct {
	AppEnv string
	DBUrl  string
	Port   string
	Redis  RedisConfig
}

func Load() *Config {
	_ = godotenv.Load()
	port, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	cfg := &Config{
		AppEnv: os.Getenv("APP_ENV"),
		DBUrl:  os.Getenv("DATABASE_URL"),
		Port:   os.Getenv("APP_PORT"),
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     port,
			DB:       db,
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			TLS:      os.Getenv("REDIS_TLS") == "true",
		},
	}

	if cfg.DBUrl == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if cfg.Port == "" {
		cfg.Port = "3000"
	}

	return cfg
}
