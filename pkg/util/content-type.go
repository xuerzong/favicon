package util

import (
	"slices"
	"strings"
)

var imageCts = []string{
	"image/png",
	"image/jpeg",
	"image/jpg",
	"image/gif",
	"image/svg+xml",
	"image/vnd.microsoft.icon",
	"image/webp",
	"image/x-icon",
}

func IsImageContentType(ct string) bool {
	return slices.Contains(imageCts, ct)
}

func ContentTypeForExtension(ext string) string {
	switch strings.ToLower(ext) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".bmp":
		return "image/bmp"
	case ".avif":
		return "image/avif"
	case ".tiff":
		return "image/tiff"
	default:
		return "application/octet-stream"
	}
}

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
