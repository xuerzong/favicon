package util

import "strings"

func ExtensionForContentType(ct string) string {
	ct = strings.TrimSpace(ct)
	if idx := strings.Index(ct, ";"); idx != -1 {
		ct = ct[:idx]
	}

	switch strings.ToLower(ct) {
	case "image/png":
		return ".png"
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	case "image/x-icon", "image/vnd.microsoft.icon":
		return ".ico"
	case "image/bmp":
		return ".bmp"
	case "image/avif":
		return ".avif"
	case "image/tiff":
		return ".tiff"
	default:
		return ".png"
	}
}
