package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Address         string
	ShutdownTimeout time.Duration
	GinMode         string
	ImageSavePath   string
	AppBaseURL      string
	CacheTTL        time.Duration
	LocalCache      bool
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	shutdownTimeout := 10 * time.Second
	if v := os.Getenv("SHUTDOWN_TIMEOUT"); v != "" {
		seconds, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		shutdownTimeout = time.Duration(seconds) * time.Second
	}

	cacheTTL := 1 * time.Hour
	if v := os.Getenv("CACHE_TTL"); v != "" {
		ttl, err := time.ParseDuration(v)
		if err != nil {
			return nil, err
		}
		cacheTTL = ttl
	}

	return &Config{
		Address:         env("ADDRESS", ":8080"),
		ShutdownTimeout: shutdownTimeout,
		GinMode:         env("GIN_MODE", "debug"),
		ImageSavePath:   env("IMAGE_SAVE_PATH", "images"),
		AppBaseURL:      env("APP_BASE_URL", "http://127.0.0.1:8080"),
		CacheTTL:        cacheTTL,
		LocalCache:      envBool("LOCAL_CACHE", true),
	}, nil
}

func env(key string, defaultVal string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	return v
}

func envBool(key string, defaultVal bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return defaultVal
	}
	return b
}
