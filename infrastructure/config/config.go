package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	DatabaseURL  string
	JWTSecret    string
	JWTTTL       time.Duration
	S3Endpoint   string
	S3Bucket     string
	S3AccessKey  string
	S3SecretKey  string
	S3Region     string
	Port         string
	APIBasePath  string
	TicketPrefix string
}

// Load reads configuration from the environment (and optionally a .env file).
func Load() (Config, error) {
	_ = godotenv.Load()

	ttlStr := getEnv("JWT_TTL", "24h")
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid JWT_TTL %q: %w", ttlStr, err)
	}

	cfg := Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		JWTTTL:       ttl,
		S3Endpoint:   os.Getenv("S3_ENDPOINT"),
		S3Bucket:     os.Getenv("S3_BUCKET"),
		S3AccessKey:  os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey:  os.Getenv("S3_SECRET_KEY"),
		S3Region:     getEnv("S3_REGION", "us-east-1"),
		Port:         getEnv("PORT", "8080"),
		APIBasePath:  getEnv("API_BASE_PATH", "/v1"),
		TicketPrefix: getEnv("TICKET_PREFIX", "40% OHHR-VOL.2"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
