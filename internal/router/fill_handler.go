package router

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/avbru/image-previewer/internal/models"
	"github.com/avbru/image-previewer/internal/services"

	"github.com/rs/zerolog/log"
)

type fillHandler struct {
	cache services.CacheService
}

func newFillHandler(cache services.CacheService) fillHandler {
	return fillHandler{
		cache: cache,
	}
}

func (h fillHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpMethodNotAllowed(w, r.Method, errors.New("method not allowed"))
		return
	}

	img, err := parseURL(r.URL)
	if err != nil {
		httpBadRequest(w, "cannot parse request URL params", err)
		return
	}

	log.Info().Msgf("image requested: %dx%d %s ", img.Width, img.Height, img.URL)

	err = h.cache.GetWithContext(r.Context(), img, r.Header, w)
	if err != nil {
		httpInternalServerError(w, "internal server error", err)
	}
}

func parseURL(uri *url.URL) (models.Image, error) {
	var head string
	var img models.Image
	var err error
	// Parse desired image width
	head, uri.Path = shiftPath(uri.Path)
	if img.Width, err = strconv.Atoi(head); err != nil {
		return models.Image{}, err
	}

	// Parse desired image height
	head, uri.Path = shiftPath(uri.Path)
	if img.Height, err = strconv.Atoi(head); err != nil {
		return models.Image{}, err
	}

	// Parse desired image URL
	img.URL = "http://" + strings.TrimPrefix(uri.Path, "/")
	imgURI, err := url.ParseRequestURI(img.URL)
	if err != nil {
		return models.Image{}, err
	}

	if imgURI.Host == "" {
		return models.Image{}, errors.New("no host provided")
	}

	return img, nil
}
