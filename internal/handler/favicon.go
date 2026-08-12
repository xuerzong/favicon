package handler

import (
	"favicon/internal/cache"
	"favicon/internal/config"
	"favicon/internal/services"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type FaviconData struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

func GetFaviconByDomain(httpClient *http.Client, cfg *config.Config, faviconCache *cache.Cache, domain string) (*FaviconData, error) {
	if item, ok := faviconCache.Get(domain); ok {
		if _, err := os.Stat(path.Join(cfg.ImageSavePath, item.Name)); err == nil {
			return &FaviconData{
				URL:  item.URL,
				Name: item.Name,
			}, nil
		}
		faviconCache.Delete(domain)
	}

	faviconUrl, err := services.GetFaviconByDomain(httpClient, domain)
	if err != nil {
		faviconCache.Set(domain, "", "default.svg")
		return nil, err
	}

	filePath, err := services.DownloadFavicon(httpClient, faviconUrl, path.Join(cfg.ImageSavePath, domain))
	if err != nil {
		faviconCache.Set(domain, "", "default.svg")
		return nil, err
	}

	name := filepath.Base(filePath)
	faviconCache.Set(domain, faviconUrl, name)

	return &FaviconData{
		URL:  faviconUrl,
		Name: name,
	}, nil
}
