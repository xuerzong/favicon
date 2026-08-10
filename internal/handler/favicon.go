package handler

import (
	"favicon/internal/config"
	"favicon/internal/services"
	"favicon/internal/util"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func GetFavicon(ctx *gin.Context, httpClient *http.Client, cfg *config.Config, siteUrl string) (any, error) {
	faviconUrl, err := services.GetFavicon(httpClient, siteUrl)
	if err != nil {
		return nil, err
	}

	domain, err := util.GetDomainFromURL(siteUrl)
	if err != nil {
		return nil, err
	}

	if err := services.GetAndSaveFavicon(httpClient, faviconUrl, path.Join(cfg.ImageSavePath, domain+".png")); err != nil {
		return nil, err
	}

	return gin.H{
		"url": faviconUrl,
	}, nil
}
