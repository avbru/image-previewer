package router

import (
	"errors"
	"net/http"

	"github.com/avbru/image-previewer/internal/services"
)

type apiHandler struct {
	cache services.CacheService
}

func newAPIHandler(cache services.CacheService) apiHandler {
	return apiHandler{
		cache: cache,
	}
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpMethodNotAllowed(w, r.Method, errors.New(" method not allowed"))
		return
	}
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "stats":
		stats, err := h.cache.GetStats()
		if err != nil {
			httpInternalServerError(w, "can't get stats", err)
			return
		}
		httpJSON(w, stats)
	case "samples":
		if _, err := w.Write([]byte(html)); err != nil {
			httpInternalServerError(w, "", err)
			return
		}
	default:
		http.NotFound(w, r)
	}
}

const html = `
<!DOCTYPE HTML>
<html>
<head>
<meta charset="utf-8">
<title>Previewer sample</title>
</head>
<body>
<h1>Sample</h1>
<p><img src="http://localhost/fill/1024/504/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg"></p>
<p><img src="http://localhost/fill/50/50/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg"></p>
<p><img src="http://localhost/fill/256/126/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg"></p>
<p><img src="http://localhost/fill/333/666/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg"></p>
</body>
</html>
`
