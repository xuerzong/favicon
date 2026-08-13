package main

import (
	"context"
	"favicon/internal/cache"
	"favicon/internal/config"
	"favicon/internal/service"
	"favicon/pkg/core"
	"favicon/pkg/util"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
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

	if cfg.LocalCache {
		go service.StartCacheCleanup(cfg.ImageSavePath, cfg.CacheTTL, logger)
	}

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
		raw := false
		siteUrl := strings.TrimPrefix(r.URL.Path, "/")
		if strings.HasPrefix(siteUrl, "raw/") {
			raw = true
			siteUrl = strings.TrimPrefix(siteUrl, "raw/")
		}

		if siteUrl == "" {
			http.ServeFile(w, r, path.Join(cfg.ImageSavePath, "default.svg"))
			return
		}

		domain, err := util.GetDomainFromURL(siteUrl)
		if err != nil {
			http.ServeFile(w, r, path.Join(cfg.ImageSavePath, "default.svg"))
			return
		}

		if !raw && cfg.ImgproxyURL != "" {
			serveViaImgproxy(w, r, client, cfg, domain)
			return
		}

		serveRawFavicon(w, r, client, cfg, faviconCache, domain)
	})
}

func serveRawFavicon(w http.ResponseWriter, r *http.Request, client *http.Client, cfg *config.Config, faviconCache *cache.Cache, domain string) {
	fv := core.NewFavicon(client, domain)

	data, err, _ := sf.Do(domain, func() (any, error) {
		return service.GetFaviconByDomain(cfg, faviconCache, fv)
	})

	if err != nil {
		http.ServeFile(w, r, path.Join(cfg.ImageSavePath, "default.svg"))
		return
	}

	fdata := data.(*service.FaviconResult)

	if cfg.LocalCache {
		http.ServeFile(w, r, path.Join(cfg.ImageSavePath, fdata.Name))
		return
	}

	w.Header().Set("Content-Type", fdata.Data.ContentType)
	w.Write(fdata.Data.Data)
}

func serveViaImgproxy(w http.ResponseWriter, r *http.Request, client *http.Client, cfg *config.Config, domain string) {
	imgproxy := service.NewImgproxy(cfg.ImgproxyURL, cfg.ImgproxySourceURL+domain, parseImgproxyOpts(r))

	resp, err := client.Get(imgproxy.Build())
	if err != nil {
		http.ServeFile(w, r, path.Join(cfg.ImageSavePath, "default.svg"))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.ServeFile(w, r, path.Join(cfg.ImageSavePath, "default.svg"))
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Cache-Control", resp.Header.Get("Cache-Control"))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, resp.Body)
}

func parseImgproxyOpts(r *http.Request) *service.ImgproxyOpts {
	q := r.URL.Query()
	opts := &service.ImgproxyOpts{
		Format: q.Get("format"),
	}
	if n, err := strconv.Atoi(q.Get("size")); err == nil {
		opts.Size = n
	}
	if n, err := strconv.Atoi(q.Get("quality")); err == nil {
		opts.Quality = n
	}
	if n, err := strconv.Atoi(q.Get("rotate")); err == nil {
		opts.Rotate = n
	}
	return opts
}
