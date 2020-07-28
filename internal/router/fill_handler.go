package router

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/avbru/image-previewer/internal/cache"
	"github.com/avbru/image-previewer/internal/resizer"
	"github.com/rs/zerolog/log"
)

type fillHandler struct {
	cache cache.Cache
}

func newFillHandler(cache cache.Cache) fillHandler {
	return fillHandler{
		cache: cache,
	}
}

func (h fillHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpMethodNotAllowed(w, "method not allowed", nil)
		return
	}

	width, height, imgURI, err := parseURL(r)
	if err != nil {
		httpBadRequest(w, "cannot parse request URL params", err)
		return
	}

	log.Info().Msgf("image requested: %dx%d %s ", width, height, imgURI)
	//TODO check in cache

	client := http.Client{}
	req, err := http.NewRequestWithContext(context.TODO(), "GET", imgURI, strings.NewReader(""))
	if err != nil {
		httpInternalServerError(w, "cannot create new request to remote host", err)
		return
	}
	copyHeaders(req.Header, r.Header)

	resp, err := client.Do(req)
	if err != nil {
		httpRemoteServerError(w, "remote server unreachable", err)
		return
	}
	defer resp.Body.Close()

	img, err := resizer.Resize(width, height, resp.Body)
	if err != nil {
		httpInternalServerError(w, "cannot resize image", err)
		return
	}
	_, err = io.Copy(w, img)
	if err != nil {
		httpInternalServerError(w, "cannot deliver image", err)
		return
	}

	//TODO add to cache
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func parseURL(r *http.Request) (width int, height int, imgURL string, err error) {
	//Parse desired image width
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	width, err = strconv.Atoi(head)
	if err != nil {
		return
	}

	//Parse desired image height
	head, r.URL.Path = shiftPath(r.URL.Path)
	height, err = strconv.Atoi(head)
	if err != nil {
		return
	}

	//Parse desired image URL
	remoteURL := "http://" + strings.TrimPrefix(r.URL.Path, "/")
	imgURI, err := url.ParseRequestURI(remoteURL)
	if err != nil {
		return
	}

	if imgURI.Host == "" {
		return
	}

	return width, height, imgURI.String(), nil
}
