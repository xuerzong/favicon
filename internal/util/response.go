package util

import "slices"

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
