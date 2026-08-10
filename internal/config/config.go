package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string

	APIHost string
	APIPort string
	WebURL  string

	DatabaseURL string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3UseSSL    bool

	SessionCookieName string
	SessionTTL        time.Duration
	SessionSecure     bool
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return nil, fmt.Errorf("parse REDIS_DB: %w", err)
	}

	s3UseSSL, err := strconv.ParseBool(getEnv("S3_USE_SSL", "false"))
	if err != nil {
		return nil, fmt.Errorf("parse S3_USE_SSL: %w", err)
	}

	sessionTTLHours, err := strconv.Atoi(getEnv("SESSION_TTL_HOURS", "168"))
	if err != nil {
		return nil, fmt.Errorf("parse SESSION_TTL_HOURS: %w", err)
	}

	sessionSecure, err := strconv.ParseBool(getEnv("SESSION_SECURE", "false"))
	if err != nil {
		return nil, fmt.Errorf("parse SESSION_SECURE: %w", err)
	}

	if sessionTTLHours <= 0 {
		return nil, fmt.Errorf("SESSION_TTL_HOURS must be greater than zero")
	}

	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "development"),

		APIHost: getEnv("API_HOST", "0.0.0.0"),
		APIPort: getEnv("API_PORT", "8080"),
		WebURL:  getEnv("WEB_URL", "http://localhost:3000"),

		DatabaseURL: os.Getenv("DATABASE_URL"),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,

		S3Endpoint:  getEnv("S3_ENDPOINT", "localhost:9000"),
		S3AccessKey: os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey: os.Getenv("S3_SECRET_KEY"),
		S3Bucket:    getEnv("S3_BUCKET", "peakcloud"),
		S3UseSSL:    s3UseSSL,

		SessionCookieName: getEnv("SESSION_COOKIE_NAME", "peakcloud_session"),
		SessionTTL:        time.Duration(sessionTTLHours) * time.Hour,
		SessionSecure:     sessionSecure,
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
