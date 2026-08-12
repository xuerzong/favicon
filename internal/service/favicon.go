package service

import (
	"favicon/pkg/core"
	"favicon/pkg/util"
	"os"

	"favicon/internal/cache"
	"favicon/internal/config"
	"path"
	"path/filepath"
)

type FaviconResult struct {
	Name string
	Data *core.FaviconData
}

func saveLocalCache(cfg *config.Config, faviconCache *cache.Cache, domain string, fd *core.FaviconData) (string, error) {
	filePath := path.Join(cfg.ImageSavePath, domain) + util.ExtensionForContentType(fd.ContentType)
	if err := os.WriteFile(filePath, fd.Data, 0644); err != nil {
		faviconCache.Set(domain, "", "default.svg")
		return "", err
	}
	name := filepath.Base(filePath)
	faviconCache.Set(domain, fd.URL, name)
	return name, nil
}

func GetFaviconByDomain(cfg *config.Config, faviconCache *cache.Cache, fv *core.Favicon) (*FaviconResult, error) {
	domain := fv.GetDomain()
	if item, ok := faviconCache.Get(domain); ok {
		fp := path.Join(cfg.ImageSavePath, item.Name)
		filedata, err := os.ReadFile(fp)

		if err == nil {
			return &FaviconResult{
				Name: item.Name,
				Data: &core.FaviconData{
					URL:         item.URL,
					Data:        filedata,
					ContentType: util.ContentTypeForExtension(path.Ext(fp)),
				},
			}, nil
		}

		fd, err := fv.FetchFaviconData(item.URL)
		if err == nil {
			if cfg.LocalCache {
				name, err := saveLocalCache(cfg, faviconCache, domain, fd)
				if err != nil {
					return nil, err
				}
				return &FaviconResult{Name: name, Data: fd}, nil
			}
			return &FaviconResult{Data: fd}, nil
		}
	}

	fd, err := fv.Get()
	if err != nil {
		faviconCache.Set(domain, "", "default.svg")
		return nil, err
	}

	name := ""
	if cfg.LocalCache {
		name, err = saveLocalCache(cfg, faviconCache, domain, fd)
		if err != nil {
			return nil, err
		}
	}

	return &FaviconResult{
		Data: fd,
		Name: name,
	}, nil
}
