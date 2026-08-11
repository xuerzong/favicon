package main

import (
	"context"
	"favicon/internal/config"
	"favicon/internal/handler"
	"favicon/internal/response"
	"favicon/internal/util"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

var sf singleflight.Group

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	gin.SetMode(cfg.GinMode)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	client := &http.Client{
		Timeout: 8 * time.Second,
	}

	router.GET("/favicon", func(ctx *gin.Context) {
		siteUrl, siteUrlOk := ctx.GetQuery("url")
		if !siteUrlOk {
			response.Error(ctx, response.NewResponseError(http.StatusBadRequest, "BAD_REQUEST"))
			return
		}

		domain, err := util.GetDomainFromURL(siteUrl)
		if err != nil {
			response.Error(ctx, response.NewResponseError(http.StatusInternalServerError, err.Error()))
			return
		}

		data, err, _ := sf.Do(siteUrl, func() (any, error) {
			return handler.GetFaviconByDomain(client, cfg, domain)
		})

		if err != nil {
			response.Error(ctx, err)
			return
		}
		response.Success(ctx, data)
	})

	router.NoRoute(func(ctx *gin.Context) {
		siteUrl := ctx.Request.URL.Path[1:]
		if siteUrl == "" {
			response.Error(ctx, response.NewResponseError(http.StatusBadRequest, "BAD_REQUEST"))
			return
		}

		domain, err := util.GetDomainFromURL(siteUrl)
		if err != nil {
			response.Error(ctx, response.NewResponseError(http.StatusInternalServerError, err.Error()))
			return
		}

		data, err, _ := sf.Do(domain, func() (any, error) {
			return handler.GetFaviconByDomain(client, cfg, domain)
		})

		fdata := data.(*handler.FaviconData)

		if err != nil {
			response.Error(ctx, err)
			return
		}
		ctx.File(path.Join(cfg.ImageSavePath, fdata.Name))
	})

	router.NoMethod(func(ctx *gin.Context) {
		response.Error(ctx, response.NewResponseError(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED"))
	})

	server := New(router, cfg.Address, cfg.ShutdownTimeout, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
