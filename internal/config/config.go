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

	return &Config{
		Address:         env("ADDRESS", ":8080"),
		ShutdownTimeout: shutdownTimeout,
		GinMode:         env("GIN_MODE", "debug"),
		ImageSavePath:   env("IMAGE_SAVE_PATH", "images"),
		AppBaseURL:      env("APP_BASE_URL", "http://127.0.0.1:8080"),
	}, nil
}

func env(key string, defaultVal string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	return v
}
