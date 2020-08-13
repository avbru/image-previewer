// +build integration

package integration

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Errors(t *testing.T) {
	srv := os.Getenv("PREVIEWER_URL")
	require.NotEmpty(t, srv, "server url not provided.")

	client := http.Client{}
	for _, tt := range queries {
		tCase := tt
		t.Run(tCase.URL, func(t *testing.T) {
			req, err := http.NewRequestWithContext(context.Background(), "GET", srv+tCase.URL, nil)
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			require.Contains(t, string(body), tCase.Body)
		})
	}
}

var queries = []struct {
	URL  string
	Body string
}{
	{
		URL:  "/fill/100/100/fileserver/static/large_gopher.jpg",
		Body: "too large",
	},
	{
		URL:  "/fill/120/120/fileserver/static/1.txt",
		Body: "wrong format",
	},
	{
		URL:  "/fill/150/150/fileserver/static/no.jpg",
		Body: "wrong format",
	},
	{
		URL:  "/fill/150/150/fileserver1/static/no.jpg",
		Body: "remote",
	},
}
