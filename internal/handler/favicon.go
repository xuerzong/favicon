package handler

import (
	"favicon/internal/config"
	"favicon/internal/services"
	"net/http"
	"path"
)

type FaviconData struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

func GetFaviconByDomain(httpClient *http.Client, cfg *config.Config, domain string) (*FaviconData, error) {
	faviconUrl, err := services.GetFaviconByDomain(httpClient, domain)
	if err != nil {
		return nil, err
	}

	filename := domain + ".png"

	if err := services.DownloadFavicon(httpClient, faviconUrl, path.Join(cfg.ImageSavePath, filename)); err != nil {
		return nil, err
	}

	return &FaviconData{
		URL:  faviconUrl,
		Name: filename,
	}, nil
}
