package service

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type ImgproxyOpts struct {
	Size    int
	Quality int
	Format  string
	Rotate  int
}

type Imgproxy struct {
	baseURL string
	source  string
	opts    *ImgproxyOpts
}

func NewImgproxy(baseURL, source string, opts *ImgproxyOpts) *Imgproxy {
	return &Imgproxy{baseURL: baseURL, source: source, opts: opts}
}

func (i *Imgproxy) Build() string {
	encoded := base64.RawURLEncoding.EncodeToString([]byte(i.source))

	u := strings.TrimRight(i.baseURL, "/") + "/unsafe"
	if segments := i.opts.segments(); segments != "" {
		u += "/" + segments
	}
	return u + "/" + encoded
}

func (o *ImgproxyOpts) segments() string {
	var opts []string

	if o.Size >= 1 && o.Size <= 2500 {
		opts = append(opts, fmt.Sprintf("resizing_type:fill/width:%d/height:%d", o.Size, o.Size))
	}

	if o.Quality >= 20 && o.Quality <= 100 {
		opts = append(opts, fmt.Sprintf("quality:%d", o.Quality))
	}

	if o.Format != "" && o.Format != "origin" {
		switch o.Format {
		case "jpg", "jpeg", "png", "webp", "avif", "gif":
			format := o.Format
			if format == "jpeg" {
				format = "jpg"
			}
			opts = append(opts, "format:"+format)
		}
	}

	if o.Rotate >= 0 && o.Rotate <= 360 && o.Rotate%90 == 0 {
		opts = append(opts, fmt.Sprintf("rotate:%d", o.Rotate))
	}

	return strings.Join(opts, "/")
}
