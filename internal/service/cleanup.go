package service

import (
	"log/slog"
	"os"
	"path"
	"time"
)

func StartCacheCleanup(imageSavePath string, cacheTTL time.Duration, logger *slog.Logger) {
	cleanupImageCache(imageSavePath, cacheTTL, logger)

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		cleanupImageCache(imageSavePath, cacheTTL, logger)
	}
}

func cleanupImageCache(imageSavePath string, cacheTTL time.Duration, logger *slog.Logger) {
	entries, err := os.ReadDir(imageSavePath)
	if err != nil {
		logger.Warn("cache cleanup: failed to read image dir", "error", err)
		return
	}

	cutoff := time.Now().Add(-cacheTTL)
	deleted := 0
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "default.svg" {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(cutoff) {
			continue
		}
		if err := os.Remove(path.Join(imageSavePath, entry.Name())); err != nil {
			logger.Warn("cache cleanup: failed to remove file", "file", entry.Name(), "error", err)
			continue
		}
		deleted++
	}

	if deleted > 0 {
		logger.Info("cache cleanup", "deleted", deleted)
	}
}
