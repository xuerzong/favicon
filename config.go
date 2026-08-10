package main

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

	address := os.Getenv("ADDRESS")
	if address == "" {
		address = ":8080"
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}

	imageSavePath := os.Getenv("IMAGE_SAVE_PATH")
	if imageSavePath == "" {
		imageSavePath = "images"
	}

	return &Config{
		Address:         address,
		ShutdownTimeout: shutdownTimeout,
		GinMode:         ginMode,
		ImageSavePath:   imageSavePath,
	}, nil
}
