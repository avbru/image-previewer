package router

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
		wantStatusCode: http.StatusUnprocessableEntity,
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
	srv := httptest.NewServer(newRootHandler(nil))
	client := http.Client{}
	for _, tt := range fillTests {
		tCase := tt
		t.Run(tCase.method+" "+tCase.url, func(t *testing.T) {
			req, err := http.NewRequest(tCase.method, srv.URL+tCase.url, strings.NewReader("ff"))

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
