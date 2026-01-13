package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
	Port  string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DBUrl: os.Getenv("DATABASE_URL"),
		Port:  os.Getenv("APP_PORT"),
	}

	if cfg.DBUrl == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if cfg.Port == "" {
		cfg.Port = "3000"
	}

	return cfg
}
