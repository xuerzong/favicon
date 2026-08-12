package util

import (
	"net/url"
	"strings"
)

func EnsureHTTPS(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return input
	}
	return "https://" + input
}

func GetDomainFromURL(rawURL string) (string, error) {
	rawURL = EnsureHTTPS(rawURL)
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	host := u.Host
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}
	return host, nil
}
