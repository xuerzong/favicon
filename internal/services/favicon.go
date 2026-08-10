package services

import (
	"favicon/internal/util"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

const USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) GoFetcher/1.0"

func GetFavicon(client *http.Client, siteUrl string) (string, error) {
	siteUrl = util.EnsureHTTPS(siteUrl)

	req, err := http.NewRequest(http.MethodGet, siteUrl, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", USER_AGENT)

	resp, err := client.Do(req)
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

	fmt.Println(candidates)

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

func GetFaviconFromDuckDuckGo(bareDomain string) string {
	return fmt.Sprintf("https://icons.duckduckgo.com/ip3/%s.ico", bareDomain)
}

func GetFaviconFromGoogle(bareDomain string, size int) string {
	return fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=%d", bareDomain, size)
}

func GetAndSaveFavicon(client *http.Client, faviconUrl string, filePath string) error {
	req, err := http.NewRequest(http.MethodGet, faviconUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", USER_AGENT)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")

	if !util.IsImageContentType(ct) {
		return fmt.Errorf("Content type is: %s", ct)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
