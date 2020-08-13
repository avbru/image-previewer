package cacheservice

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/avbru/image-previewer/internal/models"

	"github.com/stretchr/testify/require"
)

func Test_Opts(t *testing.T) {
	_, err := New(WithPath(""))
	require.Error(t, err)
	_, err = New(WithFilter("unknown"))
	require.Error(t, err)
	_, err = New(WithMaxImageSize(0))
	require.Error(t, err)
	s, err := New(
		WithPath("images_test"),
		WithFilter("Box"),
		WithCacheSize(10),
		WithMaxImageSize(1024),
	)
	require.NoError(t, err)
	require.Equal(t, int64(1024), s.maxImageSize)

	_ = os.Remove("images_test")
}

func Test_Headers(t *testing.T) {
	wantHeader := http.Header{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wantHeader = r.Header
	}))
	defer ts.Close()

	service, err := New(WithPath("images_test"), WithCacheSize(10))
	require.NoError(t, err)

	img := models.Image{
		URL: ts.URL,
	}

	r := http.Header{}
	r.Add("MyHeader", "my header value") // Add new header
	r.Add("User-Agent", "my user agent") // Replace header defaults
	r.Add("Accept-Encoding", "no")

	buf := bytes.Buffer{}
	w := bufio.NewWriter(&buf)
	err = service.resizeAndCache(context.Background(), img, r, w)
	require.Error(t, err)
	require.Equal(t, wantHeader, r)
	_ = os.Remove("images_test")
}

func Test_LimitReader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("large........ body"))
	}))
	defer ts.Close()

	service, err := New(WithPath("images_test"), WithMaxImageSize(1))
	require.NoError(t, err)

	img := models.Image{
		URL: ts.URL,
	}

	buf := bytes.Buffer{}
	w := bufio.NewWriter(&buf)
	err = service.resizeAndCache(context.Background(), img, http.Header{}, w)
	require.Error(t, err)
	_ = os.Remove("images_test")
}

func TestCacheService_GetWithContext_Async(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("../../../test/images/" + r.URL.Path)
		require.NoError(t, err)
		_, _ = io.Copy(w, file)
		file.Close()
	}))
	defer ts.Close()

	// Generate unique dir to avoid test runs interference
	// pathUUID,err:=uuid.NewUUID()
	// require.NoError(t,err)
	path := "images_async"

	service, err := New(WithPath(path), WithMaxImageSize(10*1024*1024), WithCacheSize(2))
	require.NoError(t, err)
	defer clean(t, service, path)

	var wg sync.WaitGroup
	wg.Add(len(images))
	for _, img := range images {
		image := img
		go func(wg *sync.WaitGroup) {
			image.URL = ts.URL + image.URL
			buf := bytes.Buffer{}
			w := bufio.NewWriter(&buf)
			err1 := service.GetWithContext(context.Background(), image, http.Header{}, w)
			require.NoError(t, err1)
			wg.Done()
		}(&wg)
	}
	wg.Wait()

	// Only 2 files should leave in cache
	time.Sleep(time.Second * 1) // Need some time to delete files
	files, _ := ioutil.ReadDir(path)
	require.Equal(t, 2, len(files))
}

func TestCacheService_GetWithContext_Sync(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("../../../test/images/" + r.URL.Path)
		require.NoError(t, err)
		_, _ = io.Copy(w, file)
		file.Close()
	}))
	defer ts.Close()

	// pathUUID,err:=uuid.NewUUID()
	// require.NoError(t,err)
	path := "images_sync"

	service, err := New(WithPath(path), WithMaxImageSize(10*1024*1024), WithCacheSize(2))
	require.NoError(t, err)
	defer clean(t, service, path)

	for _, img := range images {
		image := img
		image.URL = ts.URL + image.URL
		buf := bytes.Buffer{}
		w := bufio.NewWriter(&buf)
		err1 := service.GetWithContext(context.Background(), image, http.Header{}, w)
		require.NoError(t, err1)
		time.Sleep(time.Second)
	}
	// Test Stats
	stats, err := service.GetStats()
	require.NoError(t, err)
	fmt.Printf("%+v", stats)
	require.Equal(t, uint64(2), stats.HitCount, "hitCount")
	require.Equal(t, uint64(4), stats.MissCount, "missCount")

	// Only 2 files should leave in disk cache with size 120 and 150
	time.Sleep(time.Second * 1) // Need some time to delete files
	files, _ := ioutil.ReadDir(path)
	require.Equal(t, 2, len(files))
	for _, f := range files {
		file, err := os.Open(path + "/" + f.Name())
		require.NoError(t, err)
		img, _, err := image.Decode(file)
		require.NoError(t, err)
		require.InDelta(t, 135, img.Bounds().Max.X, 15)
		file.Close()
	}
}

var images = []models.Image{
	{
		URL:    "/gopher1.jpg",
		Width:  100,
		Height: 100,
	},
	{
		URL:    "/gopher1.jpg",
		Width:  100,
		Height: 100,
	},
	{
		URL:    "/gopher1.jpg",
		Width:  100,
		Height: 100,
	},
	{
		URL:    "/gopher1.jpg",
		Width:  110,
		Height: 110,
	},
	{
		URL:    "/gopher1.jpg",
		Width:  120,
		Height: 120,
	},
	{
		URL:    "/gopher1.jpg",
		Width:  150,
		Height: 150,
	},
}

func clean(t *testing.T, service *CacheService, path string) {
	err := service.store.Purge(context.Background())
	require.NoError(t, err)
	_ = os.Remove(path)
}
