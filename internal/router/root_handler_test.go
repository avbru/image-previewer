package router

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var rootTests = []struct {
	method         string
	url            string
	wantStatusCode int
	wantBody       string
}{
	{
		method:         "GET",
		url:            "/",
		wantStatusCode: http.StatusNotFound,
		wantBody:       "404",
	},
	{
		method:         "GET",
		url:            "",
		wantStatusCode: http.StatusNotFound,
		wantBody:       "404",
	},
	{
		method:         "GET",
		url:            "/does-not-exist",
		wantStatusCode: http.StatusNotFound,
		wantBody:       "404",
	},
}

func Test_RootHandler(t *testing.T) {
	srv := httptest.NewServer(newRootHandler(nil))
	client := http.Client{}
	for _, tt := range rootTests {
		tCase := tt
		t.Run(tCase.url, func(t *testing.T) {
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
