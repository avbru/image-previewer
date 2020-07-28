// +build integration

package integration

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var tests = []struct {
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
	{
		method:         "GET",
		url:            "/fill/200/200/fileserver/static/gopher1.jpg",
		wantStatusCode: http.StatusOK,
	},
}

func Test_RootHandler(t *testing.T) {
	srv := os.Getenv("PREVIEWER_URL")
	require.NotEmpty(t, srv, "server url not provided.")

	client := http.Client{}
	for _, tt := range tests {
		tCase := tt
		t.Run(tCase.url, func(t *testing.T) {
			req, err := http.NewRequest(tCase.method, srv+tCase.url, nil)
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
