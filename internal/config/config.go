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

	MaxUploadSize int64

	SessionCookieName string
	SessionTTL        time.Duration
	SessionSecure     bool

	Reliability ReliabilityConfig
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

	maxUploadMB, err := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE_MB", "100"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse MAX_UPLOAD_SIZE_MB: %w", err)
	}

	if maxUploadMB <= 0 {
		return nil, fmt.Errorf("MAX_UPLOAD_SIZE_MB must be greater than zero")
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

	reliability, err := loadReliabilityConfig()
	if err != nil {
		return nil, err
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

		MaxUploadSize: maxUploadMB * 1024 * 1024,

		SessionCookieName: getEnv("SESSION_COOKIE_NAME", "peakcloud_session"),
		SessionTTL:        time.Duration(sessionTTLHours) * time.Hour,
		SessionSecure:     sessionSecure,

		Reliability: reliability,
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.S3AccessKey == "" {
		return nil, fmt.Errorf("S3_ACCESS_KEY is required")
	}

	if cfg.S3SecretKey == "" {
		return nil, fmt.Errorf("S3_SECRET_KEY is required")
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
