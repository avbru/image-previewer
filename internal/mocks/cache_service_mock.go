package mocks

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/avbru/image-previewer/internal/models"
)

type CacheServiceMock struct {
}

func (s *CacheServiceMock) GetWithContext(ctx context.Context, img models.Image, r http.Header, w io.Writer) error {
	if img.URL == "http://non-existing-site.com" {
		return errors.New("remote host unreachable")
	}
	return nil
}

func (s *CacheServiceMock) GetStats() (models.Stat, error) {
	return models.Stat{
		HitCount:  1,
		MissCount: 1,
		Size:      1,
		Volume:    1,
	}, nil
}
