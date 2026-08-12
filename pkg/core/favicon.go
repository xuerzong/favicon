package core

import (
	"encoding/base64"
	"errors"
	"favicon/pkg/util"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) GoFetcher/1.0"

type Favicon struct {
	client *http.Client
	domain string
}

func NewFavicon(client *http.Client, domain string) *Favicon {
	return &Favicon{client, domain}
}

func (f *Favicon) GetDomain() string {
	return f.domain
}

func (f *Favicon) getFaviconUrlFromSite() (string, error) {
	domain := util.EnsureHTTPS(f.domain)

	req, err := http.NewRequest(http.MethodGet, domain, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", USER_AGENT)

	resp, err := f.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Http error: %d", resp.StatusCode)
	}

	baseUrl := resp.Request.URL

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	var candidates []string
	doc.Find("link[rel*='apple-touch-icon']").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok && href != "" {
			candidates = append(candidates, href)
		}
	})

	doc.Find("link[rel*='icon']").Each(func(i int, s *goquery.Selection) {
		rel, _ := s.Attr("rel")
		if rel == "apple-touch-icon" {
			return
		}
		href, ok := s.Attr("href")
		if ok && href != "" {
			candidates = append(candidates, href)
		}
	})

	for _, href := range candidates {
		abs, err := baseUrl.Parse(href)
		if err == nil {
			return abs.String(), nil
		}
	}

	fallback, err := baseUrl.Parse("/favicon.ico")

	if err != nil {
		return "", err
	}

	return fallback.String(), nil
}

func (f *Favicon) getFaviconUrlFromDuckDuckGo() string {
	return fmt.Sprintf("https://icons.duckduckgo.com/ip3/%s.ico", f.domain)
}

func (f *Favicon) getFaviconUrlFromGoogle(size int) string {
	return fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=%d", f.domain, size)
}

func (f *Favicon) getFaviconUrls() []string {
	urls := []string{}
	if url, err := f.getFaviconUrlFromSite(); err == nil && url != "" {
		urls = append(urls, url)
	}
	urls = append(urls, f.getFaviconUrlFromGoogle(64),
		f.getFaviconUrlFromDuckDuckGo())
	return urls
}

type FaviconData struct {
	Data        []byte
	URL         string
	ContentType string
}

func (f *Favicon) Get() (*FaviconData, error) {
	urls := f.getFaviconUrls()
	for _, u := range urls {
		fd, err := f.FetchFaviconData(u)
		if err != nil {
			continue
		}
		return fd, nil
	}
	return nil, fmt.Errorf("no favicon found for %s", f.domain)
}

func (f *Favicon) FetchFaviconData(faviconUrl string) (*FaviconData, error) {
	if faviconUrl == "" {
		return nil, errors.New("Favicon url is empty")
	}

	if strings.HasPrefix(faviconUrl, "data:") {
		data, err := f.decodeDataURI(faviconUrl)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	req, err := http.NewRequest(http.MethodGet, faviconUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", USER_AGENT)

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !util.IsImageContentType(ct) {
		return nil, fmt.Errorf("content type is: %s", ct)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("empty body for %s", faviconUrl)
	}
	if detected := http.DetectContentType(data); !util.IsImageContentType(detected) {
		return nil, fmt.Errorf("content %s is not an image", detected)
	}

	return &FaviconData{
		URL:         faviconUrl,
		Data:        data,
		ContentType: ct,
	}, nil
}

func (f *Favicon) decodeDataURI(dataURI string) (*FaviconData, error) {
	header, payload, ok := strings.Cut(dataURI, ",")
	if !ok {
		return nil, fmt.Errorf("invalid data uri")
	}

	ct := strings.TrimPrefix(header, "data:")
	if idx := strings.Index(ct, ";"); idx != -1 {
		ct = ct[:idx]
	}

	var data = []byte{}

	if strings.HasSuffix(header, ";base64") {
		b64, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return nil, err
		}
		data = b64
	} else {
		decoded, err := url.QueryUnescape(payload)
		if err != nil {
			return nil, err
		}
		data = []byte(decoded)
	}

	return &FaviconData{
		URL:         dataURI,
		Data:        data,
		ContentType: ct,
	}, nil
}
