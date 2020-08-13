package services

import (
	"context"
	"io"
	"net/http"

	"github.com/avbru/image-previewer/internal/models"
)

type CacheService interface {
	GetWithContext(ctx context.Context, img models.Image, r http.Header, w io.Writer) error
	GetStats() (models.Stat, error)
}
