package main

import (
	"net/http"
	"os"

	// TODO _ "github.com/avbru/image-previewer/internal/cache"
	"github.com/avbru/image-previewer/internal/router"
	"github.com/rs/zerolog/log"
)

func main() {
	handler := router.New(nil) // TODO provide cache to router

	if err := http.ListenAndServe(":80", handler.RootHandler()); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
}
