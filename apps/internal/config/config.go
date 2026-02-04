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

type StorageConfig struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type Config struct {
	AppEnv  string
	DBUrl   string
	Port    string
	Redis   RedisConfig
	Storage StorageConfig
}

func Load() *Config {
	_ = godotenv.Load()

	redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	cfg := &Config{
		AppEnv: os.Getenv("APP_ENV"),
		DBUrl:  os.Getenv("DATABASE_URL"),
		Port:   os.Getenv("APP_PORT"),

		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     redisPort,
			DB:       redisDB,
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			TLS:      os.Getenv("REDIS_TLS") == "true",
		},

		Storage: StorageConfig{
			Endpoint:  os.Getenv("STORAGE_ENDPOINT"),
			AccessKey: os.Getenv("STORAGE_ACCESS_KEY"),
			SecretKey: os.Getenv("STORAGE_SECRET_KEY"),
			Bucket:    os.Getenv("STORAGE_BUCKET"),
			UseSSL:    os.Getenv("STORAGE_USE_SSL") == "true",
		},
	}

	if cfg.DBUrl == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if cfg.Port == "" {
		cfg.Port = "3000"
	}

	if cfg.Storage.Endpoint == "" {
		log.Fatal("STORAGE_ENDPOINT is required")
	}

	if cfg.Storage.Bucket == "" {
		log.Fatal("STORAGE_BUCKET is required")
	}

	return cfg
}
