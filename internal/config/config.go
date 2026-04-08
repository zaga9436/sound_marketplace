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
	AppEnv                string
	AppPort               string
	PostgresHost          string
	PostgresPort          string
	PostgresDB            string
	PostgresUser          string
	PostgresPassword      string
	PostgresSSLMode       string
	RedisHost             string
	RedisPort             string
	RedisPassword         string
	JWTSecret             string
	JWTTTL                time.Duration
	S3Endpoint            string
	S3Region              string
	S3Bucket              string
	S3AccessKey           string
	S3SecretKey           string
	S3UseSSL              bool
	SignedURLTTL          time.Duration
	YooKassaShopID        string
	YooKassaSecretKey     string
	YooKassaReturnURL     string
	MaxUploadSize         int64
	AllowedAudioFormats   []string
	AllowedImageFormats   []string
	WSReadBufferSize      int
	WSWriteBufferSize     int
	PreviewDurationSecond int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	required, err := requiredEnv(
		"APP_PORT",
		"POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_DB", "POSTGRES_USER", "POSTGRES_PASSWORD",
		"REDIS_HOST", "REDIS_PORT",
		"JWT_SECRET", "JWT_TTL",
		"S3_ENDPOINT", "S3_REGION", "S3_BUCKET", "S3_ACCESS_KEY", "S3_SECRET_KEY", "S3_USE_SSL", "SIGNED_URL_TTL",
		"YOOKASSA_SHOP_ID", "YOOKASSA_SECRET_KEY", "YOOKASSA_RETURN_URL",
		"MAX_UPLOAD_SIZE", "ALLOWED_AUDIO_FORMATS", "ALLOWED_IMAGE_FORMATS",
		"WS_READ_BUFFER_SIZE", "WS_WRITE_BUFFER_SIZE", "PREVIEW_DURATION_SECONDS",
	)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppEnv:           getWithDefault("APP_ENV", "development"),
		AppPort:          required["APP_PORT"],
		PostgresHost:     required["POSTGRES_HOST"],
		PostgresPort:     required["POSTGRES_PORT"],
		PostgresDB:       required["POSTGRES_DB"],
		PostgresUser:     required["POSTGRES_USER"],
		PostgresPassword: required["POSTGRES_PASSWORD"],
		PostgresSSLMode:  getWithDefault("POSTGRES_SSLMODE", "disable"),
		RedisHost:        required["REDIS_HOST"],
		RedisPort:        required["REDIS_PORT"],
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		JWTSecret:        required["JWT_SECRET"],
		S3Endpoint:       required["S3_ENDPOINT"],
		S3Region:         required["S3_REGION"],
		S3Bucket:         required["S3_BUCKET"],
		S3AccessKey:      required["S3_ACCESS_KEY"],
		S3SecretKey:      required["S3_SECRET_KEY"],
		YooKassaShopID:   required["YOOKASSA_SHOP_ID"],
		YooKassaSecretKey: required["YOOKASSA_SECRET_KEY"],
		YooKassaReturnURL: required["YOOKASSA_RETURN_URL"],
		AllowedAudioFormats: parseCSV(required["ALLOWED_AUDIO_FORMATS"]),
		AllowedImageFormats: parseCSV(required["ALLOWED_IMAGE_FORMATS"]),
	}

	if cfg.JWTTTL, err = time.ParseDuration(required["JWT_TTL"]); err != nil {
		return nil, fmt.Errorf("parse JWT_TTL: %w", err)
	}
	if cfg.SignedURLTTL, err = time.ParseDuration(required["SIGNED_URL_TTL"]); err != nil {
		return nil, fmt.Errorf("parse SIGNED_URL_TTL: %w", err)
	}
	if cfg.MaxUploadSize, err = strconv.ParseInt(required["MAX_UPLOAD_SIZE"], 10, 64); err != nil {
		return nil, fmt.Errorf("parse MAX_UPLOAD_SIZE: %w", err)
	}
	if cfg.WSReadBufferSize, err = strconv.Atoi(required["WS_READ_BUFFER_SIZE"]); err != nil {
		return nil, fmt.Errorf("parse WS_READ_BUFFER_SIZE: %w", err)
	}
	if cfg.WSWriteBufferSize, err = strconv.Atoi(required["WS_WRITE_BUFFER_SIZE"]); err != nil {
		return nil, fmt.Errorf("parse WS_WRITE_BUFFER_SIZE: %w", err)
	}
	if cfg.PreviewDurationSecond, err = strconv.Atoi(required["PREVIEW_DURATION_SECONDS"]); err != nil {
		return nil, fmt.Errorf("parse PREVIEW_DURATION_SECONDS: %w", err)
	}
	if cfg.S3UseSSL, err = strconv.ParseBool(required["S3_USE_SSL"]); err != nil {
		return nil, fmt.Errorf("parse S3_USE_SSL: %w", err)
	}

	return cfg, nil
}

func requiredEnv(keys ...string) (map[string]string, error) {
	result := make(map[string]string, len(keys))
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			return nil, fmt.Errorf("required env missing: %s", key)
		}
		result[key] = value
	}
	return result, nil
}

func getWithDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func parseCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
