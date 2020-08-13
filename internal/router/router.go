package router

import (
	"net/http"

	"github.com/avbru/image-previewer/internal/services"
)

type Router struct {
	rootHandler rootHandler
}

func New(cache services.CacheService) *Router {
	return &Router{
		rootHandler: newRootHandler(cache),
	}
}

func (r *Router) RootHandler() http.Handler {
	return r.rootHandler
}
