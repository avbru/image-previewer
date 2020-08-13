package router

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/avbru/image-previewer/internal/mocks"
	"github.com/stretchr/testify/require"
)

var apiTests = []struct {
	method         string
	url            string
	wantStatusCode int
	wantBody       string
}{
	{
		method:         "POST",
		url:            "/api",
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
		method:         "POST",
		url:            "/api/samples",
		wantStatusCode: http.StatusMethodNotAllowed,
		wantBody:       "",
	},
	{
		method:         "GET",
		url:            "/api/samples",
		wantStatusCode: http.StatusOK,
		wantBody:       "<img ",
	},
	{
		method:         "POST",
		url:            "/api/stats",
		wantStatusCode: http.StatusMethodNotAllowed,
		wantBody:       "",
	},
	{
		method:         "GET",
		url:            "/api/stats",
		wantStatusCode: http.StatusOK,
		wantBody:       "HitCount",
	},
}

func Test_APIHandler(t *testing.T) {
	cacheService := mocks.CacheServiceMock{}

	srv := httptest.NewServer(newRootHandler(&cacheService))

	client := http.Client{}
	for _, tt := range apiTests {
		tCase := tt
		t.Run(tCase.url+" "+tCase.method, func(t *testing.T) {
			req, err := http.NewRequestWithContext(context.Background(), tCase.method, srv.URL+tCase.url, nil)

			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			require.NoError(t, err)
			require.Equal(t, tCase.wantStatusCode, resp.StatusCode)
			require.Contains(t, string(body), tCase.wantBody)
		})
	}
}
