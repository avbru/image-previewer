// +build integration

package integration

import (
	"image"
	_ "image/jpeg"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/avbru/image-previewer/internal/models"

	"github.com/stretchr/testify/require"
)

func Test_Images_Sync(t *testing.T) {
	srv := os.Getenv("PREVIEWER_URL")
	require.NotEmpty(t, srv, "server url not provided.")

	client := http.Client{}
	for _, tt := range images {
		tCase := tt
		t.Run(tCase.URL, func(t *testing.T) {
			req, err := http.NewRequest("GET", srv+tCase.URL, nil)
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode)
			img, _, err := image.Decode(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tCase.Width, img.Bounds().Max.Y)
			require.Equal(t, tCase.Height, img.Bounds().Max.Y)
		})
	}
}

func TestCacheService_GetWithContext_Async(t *testing.T) {
	srv := os.Getenv("PREVIEWER_URL")
	require.NotEmpty(t, srv, "server url not provided.")
	client := http.Client{}

	var wg sync.WaitGroup
	wg.Add(len(images))
	for _, tt := range images {
		tCase := tt
		go func(wg *sync.WaitGroup) {
			req, err := http.NewRequest("GET", srv+tCase.URL, nil)
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			img, _, err := image.Decode(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tCase.Width, img.Bounds().Max.Y)
			require.Equal(t, tCase.Height, img.Bounds().Max.Y)
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}

var images = []models.Image{
	{
		URL:    "/fill/100/100/fileserver/static/gopher1.jpg",
		Width:  100,
		Height: 100,
	},
	{
		URL:    "/fill/120/120/fileserver/static/gopher1.jpg",
		Width:  120,
		Height: 120,
	},
	{
		URL:    "/fill/150/150/fileserver/static/gopher1.jpg",
		Width:  150,
		Height: 150,
	},
	{
		URL:    "/fill/200/200/fileserver/static/subdir/gopher2.jpg",
		Width:  200,
		Height: 200,
	},
}
