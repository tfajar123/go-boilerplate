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
	PublicURL string
}

type MailConfig struct {
	Host             string
	Port             int
	Username         string
	Password         string
	FromEmail        string
	FromName         string
	ResetPasswordURL string
}

type Config struct {
	AppEnv  string
	DBUrl   string
	Port    string
	Redis   RedisConfig
	Storage StorageConfig
	Mail    MailConfig
}

func Load() *Config {
	_ = godotenv.Load()

	redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	mailPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

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
			Region:    os.Getenv("STORAGE_REGION"),
			AccessKey: os.Getenv("STORAGE_ACCESS_KEY"),
			SecretKey: os.Getenv("STORAGE_SECRET_KEY"),
			Bucket:    os.Getenv("STORAGE_BUCKET"),
			UseSSL:    os.Getenv("STORAGE_USE_SSL") == "true",
			PublicURL: os.Getenv("STORAGE_PUBLIC_URL"),
		},

		Mail: MailConfig{
			Host:             os.Getenv("MAIL_HOST"),
			Port:             mailPort,
			Username:         os.Getenv("MAIL_USERNAME"),
			Password:         os.Getenv("MAIL_PASSWORD"),
			FromEmail:        os.Getenv("MAIL_FROM_EMAIL"),
			FromName:         os.Getenv("MAIL_FROM_NAME"),
			ResetPasswordURL: os.Getenv("MAIL_RESET_PASSWORD_URL"),
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
