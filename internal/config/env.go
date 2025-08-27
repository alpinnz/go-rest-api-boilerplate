package config

import (
	"log"
	"os"
	"strconv"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/infrastructure/databases"
	"github.com/joho/godotenv"
)

type Env struct {
	App      App
	Server   Server
	Postgres databases.PostgresConfig
	Auth     Auth
	SMTP     SMTP
}

type App struct {
	Env           string
	DefaultLocale string
}

type Server struct {
	Address string
}

type Auth struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	PasswordSecret     string
}

type SMTP struct {
	Username  string
	Password  string
	Host      string
	Port      int
	FromEmail string
}

func NewEnv() Env {
	// Load .env file if present
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found or failed to load:", err)
	} else {
		log.Println("✅ .env file loaded successfully")
	}

	return Env{
		App: App{
			Env:           getEnv("APP_ENV", "development"),
			DefaultLocale: getEnv("APP_DEFAULT_LOCALE", "id"),
		},
		Server: Server{
			Address: getEnv("SERVER_ADDRESS", ":3000"),
		},
		Postgres: databases.PostgresConfig{
			Host: getEnv("DB_PG_HOST", "localhost"),
			Port: getEnvAsInt("DB_PG_PORT", 5432),
			User: getEnv("DB_PG_USER", ""),
			Pass: getEnv("DB_PG_PASS", ""),
			Name: getEnv("DB_PG_NAME", ""),
		},
		Auth: Auth{
			AccessTokenSecret:  getEnv("AUTH_ACCESS_TOKEN_SECRET", ""),
			RefreshTokenSecret: getEnv("AUTH_REFRESH_TOKEN_SECRET", ""),
			PasswordSecret:     getEnv("AUTH_PASSWORD_SECRET", ""),
		},
		SMTP: SMTP{
			Username:  getEnv("SMTP_USERNAME", ""),
			Password:  getEnv("SMTP_PASSWORD", ""),
			Host:      getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:      getEnvAsInt("SMTP_PORT", 587),
			FromEmail: getEnv("SMTP_FROM_EMAIL", ""),
		},
	}
}

func getEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return val
	}
	return fallback
}
