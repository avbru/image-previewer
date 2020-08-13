package cacheservice

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/disintegration/imaging"

	"github.com/avbru/image-previewer/internal/cache/lru"
	"github.com/avbru/image-previewer/internal/storages/filestorage"

	"github.com/avbru/image-previewer/internal/models"
	"github.com/rs/zerolog/log"

	"github.com/avbru/image-previewer/internal/cache"
	"github.com/avbru/image-previewer/internal/resizer"
	"github.com/avbru/image-previewer/internal/storages"
)

type CacheService struct {
	sync.RWMutex

	hitCount  uint64
	missCount uint64

	maxImageSize int64
	filter       imaging.ResampleFilter

	cache cache.Cache
	store storages.FileStorage
}

func New(opts ...Option) (*CacheService, error) {
	var options options

	for _, o := range opts {
		if err := o.apply(&options); err != nil {
			return nil, err
		}
	}

	store, err := filestorage.New(options.path)
	if err != nil {
		return nil, err
	}

	cacheService := CacheService{
		maxImageSize: options.maxImageSize,
		cache:        lru.New(options.cacheSize),
		store:        store,
		// inProcess: make(map[string]*sync.RWMutex),
	}

	return &cacheService, nil
}

//GetStats returns current usage statistics
func (s *CacheService) GetStats() (models.Stat, error) {
	var stat models.Stat
	stat.HitCount = atomic.LoadUint64(&s.hitCount)
	stat.MissCount = atomic.LoadUint64(&s.missCount)

	stat.Size = s.cache.Len()

	var err error
	stat.Volume, err = s.store.Volume(context.Background())
	if err != nil {
		return models.Stat{}, err
	}
	return stat, nil
}

func (s *CacheService) GetWithContext(ctx context.Context, img models.Image, r http.Header, w io.Writer) error {
	fileName, ok := s.cache.Get(img.FullURL())
	if ok {
		err := s.store.ReadFile(ctx, fileName.(string), w)
		if err != nil {
			return fmt.Errorf("can't read cached image: %w", err)
		}

		atomic.AddUint64(&s.hitCount, 1)
		return nil
	}

	s.Lock()
	defer s.Unlock()
	err := s.resizeAndCache(ctx, img, r, w)
	if err != nil {
		return fmt.Errorf("can't cache image: %w", err)
	}

	atomic.AddUint64(&s.missCount, 1)
	return nil
}

func (s *CacheService) resizeAndCache(ctx context.Context, img models.Image, r http.Header, w io.Writer) error {
	client := http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", img.URL, nil)
	if err != nil {
		return fmt.Errorf("resize can't create request to remote host: %w", err)
	}

	copyHeaders(req.Header, r)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("resize can't do request to remote host: %w", err)
	}

	defer resp.Body.Close()

	originalImg := io.LimitReader(resp.Body, s.maxImageSize)

	image, err := resizer.Resize(img.Width, img.Height, s.filter, originalImg)
	if err != nil {
		return fmt.Errorf("remote image has wrong format or too large: %w", err)
	}

	// Serve image, while making copy to save later to cache
	var buf bytes.Buffer
	multiWriter := io.MultiWriter(&buf, w)
	if _, err := io.Copy(multiWriter, image); err != nil {
		return fmt.Errorf("resize can't serve resised image: %w", err)
	}
	s.addToCache(ctx, img, &buf)
	return nil
}

func (s *CacheService) addToCache(ctx context.Context, img models.Image, buf io.Reader) {
	err := s.store.SaveFile(ctx, img.FileName(), buf)
	if err != nil {
		log.Err(err)
		return
	}

	// Once file saved. Add to cache and
	// remove file evicted from LRU Cache
	removed, _ := s.cache.Set(img.FullURL(), img.FileName())
	if removed != nil {
		log.Info().Msg("file removed from cache: " + removed.(string))
		s.store.DeleteFile(ctx, removed.(string))
	}
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
