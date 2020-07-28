package router

import (
	"net/http"

	"github.com/avbru/image-previewer/internal/cache"
)

type apiHandler struct {
	cache cache.Cache
}

func newAPIHandler(cache cache.Cache) apiHandler {
	return apiHandler{
		cache: cache,
	}
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "stats":
		// TODO
	case "samples":
		if _, err := w.Write([]byte(html)); err != nil {
			httpInternalServerError(w, "", err)
			return
		}
	default:
		http.NotFound(w, r)
	}
}

// TODO Refactor to html.Template.
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
