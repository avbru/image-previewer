package router

import (
	"net/http"

	"github.com/avbru/image-previewer/internal/cache"
)

type Router struct {
	rootHandler rootHandler
}

func New(cache cache.Cache) *Router {
	return &Router{
		rootHandler: newRootHandler(cache),
	}
}

func (r *Router) RootHandler() http.Handler {
	return r.rootHandler
}
