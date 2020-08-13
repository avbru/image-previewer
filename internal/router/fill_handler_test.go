package router

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/avbru/image-previewer/internal/mocks"
	"github.com/avbru/image-previewer/internal/models"

	"github.com/stretchr/testify/require"
)

var fillTests = []struct {
	method         string
	url            string
	wantStatusCode int
	wantBody       string
}{
	{
		method:         "POST",
		url:            "/fill",
		wantStatusCode: http.StatusMethodNotAllowed,
		wantBody:       "405",
	},
	{
		method:         "DELETE",
		url:            "/fill",
		wantStatusCode: http.StatusMethodNotAllowed,
		wantBody:       "405",
	},
	{
		method:         "GET",
		url:            "/fill",
		wantStatusCode: http.StatusBadRequest,
		wantBody:       "400",
	},
	{
		method:         "GET",
		url:            "/fill/does-not-exist",
		wantStatusCode: http.StatusBadRequest,
		wantBody:       "400",
	},
	{
		method:         "GET",
		url:            "/fill/200",
		wantStatusCode: http.StatusBadRequest,
		wantBody:       "cannot parse request",
	},
	{
		method:         "GET",
		url:            "/fill/200/height",
		wantStatusCode: http.StatusBadRequest,
		wantBody:       "cannot parse request",
	},
	{
		method:         "GET",
		url:            "/fill/200/300/www.:wikipedia",
		wantStatusCode: http.StatusBadRequest,
		wantBody:       "cannot parse request",
	},
	{
		method:         "GET",
		url:            "/fill/200/300/non-existing-site.com",
		wantStatusCode: http.StatusInternalServerError,
		wantBody:       "",
	},
	{
		method:         "GET",
		url:            "/fill/200/300/myaquarium.top/img/banner.jpg",
		wantStatusCode: http.StatusOK,
		wantBody:       "",
	},
}

func Test_FillHandler(t *testing.T) {
	cacheService := mocks.CacheServiceMock{}

	srv := httptest.NewServer(newRootHandler(&cacheService))

	client := http.Client{}
	for _, tt := range fillTests {
		tCase := tt
		t.Run(tCase.method+" "+tCase.url, func(t *testing.T) {
			req, err := http.NewRequestWithContext(context.Background(), tCase.method, srv.URL+tCase.url, strings.NewReader("ff"))

			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tCase.wantStatusCode, resp.StatusCode)
			require.Contains(t, string(body), tCase.wantBody)
		})
	}
}

var urlTests = []struct {
	url     string
	wantErr bool
	wantImg models.Image
}{
	{
		url:     "",
		wantErr: true,
	},
	{
		url:     "100/dd/dd",
		wantErr: true,
	},
	{
		url:     "100/100/?",
		wantErr: true,
	},
	{
		url:     "100/100/yandex.ru/logo.jpg",
		wantErr: false,
		wantImg: models.Image{
			Width:  100,
			Height: 100,
			URL:    "http://yandex.ru/logo.jpg",
		},
	},
}

func Test_ParseURL(t *testing.T) {
	for _, tt := range urlTests {
		tCase := tt
		t.Run(tCase.url, func(t *testing.T) {
			uri := url.URL{
				Path: tCase.url,
			}
			img, err := parseURL(&uri)
			if tCase.wantErr {
				require.Error(t, err)
			}
			require.Equal(t, tCase.wantImg, img)
		})
	}
}
