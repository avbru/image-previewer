package router

import (
	"net/http"

	"github.com/avbru/image-previewer/internal/services"
)

type rootHandler struct {
	fillHandler fillHandler
	apiHandler  apiHandler
}

func newRootHandler(cache services.CacheService) rootHandler {
	return rootHandler{
		fillHandler: newFillHandler(cache),
		apiHandler:  newAPIHandler(cache),
	}
}

func (h rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	switch head {
	case "api":
		h.apiHandler.ServeHTTP(w, r)
	case "fill":
		h.fillHandler.ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}
