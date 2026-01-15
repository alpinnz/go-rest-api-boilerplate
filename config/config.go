package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App         AppConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	JWT         JWTConfig
	Server      ServerConfig
	CORS        CORSConfig
	RateLimiter RateLimiterConfig
	Logger      LoggerConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessTokenSecret      string
	AccessTokenExpiration  time.Duration
	RefreshTokenSecret     string
	RefreshTokenExpiration time.Duration
}

type ServerConfig struct {
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

type RateLimiterConfig struct {
	Enabled       bool
	RequestsLimit int
	WindowMinutes int
}

type LoggerConfig struct {
	Level  string
	Pretty bool
}

func Load() (*Config, error) {
	// Load .env file if exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "go-rest-api-boilerplate"),
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASS", "postgres"),
			Name:            getEnv("DB_NAME", "go-rest-api-boilerplate"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: parseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m")),
			ConnMaxIdleTime: parseDuration(getEnv("DB_CONN_MAX_IDLE_TIME", "10m")),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessTokenSecret:      getEnv("ACCESS_TOKEN_SECRET", "change-this-access-token-secret-key"),
			AccessTokenExpiration:  parseDuration(getEnv("ACCESS_TOKEN_EXPIRATION", "15m")),
			RefreshTokenSecret:     getEnv("REFRESH_TOKEN_SECRET", "change-this-refresh-token-secret-key"),
			RefreshTokenExpiration: parseDuration(getEnv("REFRESH_TOKEN_EXPIRATION", "168h")),
		},
		Server: ServerConfig{
			ReadTimeout:     parseDuration(getEnv("READ_TIMEOUT", "10s")),
			WriteTimeout:    parseDuration(getEnv("WRITE_TIMEOUT", "10s")),
			ShutdownTimeout: parseDuration(getEnv("SHUTDOWN_TIMEOUT", "5s")),
		},
		CORS: CORSConfig{
			AllowedOrigins:   parseStringSlice(getEnv("CORS_ALLOWED_ORIGINS", "*")),
			AllowedMethods:   parseStringSlice(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,PATCH,OPTIONS")),
			AllowedHeaders:   parseStringSlice(getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization,Accept,Origin,X-Requested-With")),
			AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           getEnvAsInt("CORS_MAX_AGE", 3600),
		},
		RateLimiter: RateLimiterConfig{
			Enabled:       getEnvAsBool("RATE_LIMITER_ENABLED", true),
			RequestsLimit: getEnvAsInt("RATE_LIMITER_REQUESTS", 100),
			WindowMinutes: getEnvAsInt("RATE_LIMITER_WINDOW_MINUTES", 1),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Pretty: getEnvAsBool("LOG_PRETTY", false),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	// App validation
	if c.App.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}
	if c.App.Env != "development" && c.App.Env != "staging" && c.App.Env != "production" {
		return fmt.Errorf("APP_ENV must be one of: development, staging, production")
	}

	// Database validation
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.Database.MaxOpenConns < 1 {
		return fmt.Errorf("DB_MAX_OPEN_CONNS must be at least 1")
	}
	if c.Database.MaxIdleConns < 1 {
		return fmt.Errorf("DB_MAX_IDLE_CONNS must be at least 1")
	}
	if c.Database.MaxIdleConns > c.Database.MaxOpenConns {
		return fmt.Errorf("DB_MAX_IDLE_CONNS cannot be greater than DB_MAX_OPEN_CONNS")
	}

	// Redis validation
	if c.Redis.Host == "" {
		return fmt.Errorf("REDIS_HOST is required")
	}

	// JWT validation
	if c.JWT.AccessTokenSecret == "" {
		return fmt.Errorf("ACCESS_TOKEN_SECRET is required")
	}
	if c.JWT.AccessTokenSecret == "change-this-access-token-secret-key" && c.App.Env == "production" {
		return fmt.Errorf("ACCESS_TOKEN_SECRET must be changed in production environment")
	}
	if len(c.JWT.AccessTokenSecret) < 32 {
		fmt.Println("Warning: ACCESS_TOKEN_SECRET should be at least 32 characters for better security")
	}
	if c.JWT.AccessTokenExpiration == 0 {
		return fmt.Errorf("ACCESS_TOKEN_EXPIRATION is required")
	}

	if c.JWT.RefreshTokenSecret == "" {
		return fmt.Errorf("REFRESH_TOKEN_SECRET is required")
	}
	if c.JWT.RefreshTokenSecret == "change-this-refresh-token-secret-key" && c.App.Env == "production" {
		return fmt.Errorf("REFRESH_TOKEN_SECRET must be changed in production environment")
	}
	if len(c.JWT.RefreshTokenSecret) < 32 {
		fmt.Println("Warning: REFRESH_TOKEN_SECRET should be at least 32 characters for better security")
	}
	if c.JWT.RefreshTokenExpiration == 0 {
		return fmt.Errorf("REFRESH_TOKEN_EXPIRATION is required")
	}

	// Server validation
	if c.Server.ReadTimeout == 0 {
		return fmt.Errorf("READ_TIMEOUT is required")
	}
	if c.Server.WriteTimeout == 0 {
		return fmt.Errorf("WRITE_TIMEOUT is required")
	}
	if c.Server.ShutdownTimeout == 0 {
		return fmt.Errorf("SHUTDOWN_TIMEOUT is required")
	}

	return nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if value == "true" || value == "1" {
			return true
		}
		if value == "false" || value == "0" {
			return false
		}
	}
	return defaultValue
}

func parseStringSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
