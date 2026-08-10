package main

import (
	"context"
	"favicon/internal/config"
	"favicon/internal/handler"
	"favicon/internal/response"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

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

		data, err := handler.GetFavicon(ctx, client, cfg, siteUrl)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		response.Success(ctx, data)
	})

	server := New(router, cfg.Address, cfg.ShutdownTimeout, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
