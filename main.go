package main

import (
	"context"
	"favicon/internal/cache"
	"favicon/internal/config"
	"favicon/internal/service"
	"favicon/internal/util"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sync/singleflight"
)

var sf singleflight.Group

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client := &http.Client{
		Timeout: 8 * time.Second,
	}

	faviconCache := cache.New(cfg.CacheTTL)

	server := New(withLogging(logger, faviconHandler(cfg, client, faviconCache)), cfg.Address, cfg.ShutdownTimeout, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

func withLogging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}

func faviconHandler(cfg *config.Config, client *http.Client, faviconCache *cache.Cache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siteUrl := strings.TrimPrefix(r.URL.Path, "/")
		if siteUrl == "" {
			http.ServeFile(w, r, path.Join(cfg.ImageSavePath, "default.svg"))
			return
		}

		domain, err := util.GetDomainFromURL(siteUrl)
		if err != nil {
			http.ServeFile(w, r, path.Join(cfg.ImageSavePath, "default.svg"))
			return
		}

		data, err, _ := sf.Do(domain, func() (any, error) {
			return service.GetFaviconByDomain(client, cfg, faviconCache, domain)
		})

		if err != nil {
			http.ServeFile(w, r, path.Join(cfg.ImageSavePath, "default.svg"))
			return
		}

		fdata := data.(*service.FaviconData)
		http.ServeFile(w, r, path.Join(cfg.ImageSavePath, fdata.Name))
	})
}
