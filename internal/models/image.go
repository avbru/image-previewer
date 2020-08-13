package models

import (
	"crypto/md5" //nolint:gosec
	"fmt"
)

type Image struct {
	Width  int
	Height int
	URL    string
}

func (img *Image) FileName() string {
	return fmt.Sprintf("%x.jpg", md5.Sum([]byte(img.FullURL()))) //nolint:gosec
}

func (img *Image) FullURL() string {
	return fmt.Sprintf("%d/%d/%s", img.Width, img.Height, img.URL)
}
