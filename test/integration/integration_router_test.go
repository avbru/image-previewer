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
	// Api handler
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
	// Fill handler
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

func Test_Router(t *testing.T) {
	srv := os.Getenv("PREVIEWER_URL")
	require.NotEmpty(t, srv, "server url not provided.")

	client := http.Client{}
	for _, tt := range tests {
		tCase := tt
		t.Run(tCase.url, func(t *testing.T) {
			req, err := http.NewRequest(tCase.method, srv+tCase.url, nil)
			require.NoError(t, err)

			resp, err := client.Do(req)

			//if err != nil {
			//	if !strings.Contains(err.Error(),"EOF") {
			//		return
			//	}
			//}
			require.NoError(t, err)
			require.Equal(t, tCase.wantStatusCode, resp.StatusCode)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				println("err read:", err.Error())
			}
			resp.Body.Close()
			require.Contains(t, string(body), tCase.wantBody)
		})
	}
}
