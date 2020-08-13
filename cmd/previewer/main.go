package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/avbru/image-previewer/internal/config"
	"github.com/avbru/image-previewer/internal/router"
	"github.com/avbru/image-previewer/internal/services/cacheservice"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("config: %+v", cfg)

	cacheService, err := cacheservice.New(
		cacheservice.WithCacheSize(cfg.CacheSize),
		cacheservice.WithMaxImageSize(cfg.MaxImageSize),
		cacheservice.WithPath(cfg.ImageDir),
		cacheservice.WithFilter(cfg.Filter),
	)
	if err != nil {
		log.Fatal().Err(err)
	}

	handler := router.New(cacheService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler.RootHandler(),
	}

	if err := runSrv(srv); err != nil {
		log.Err(err)
	}
}

func runSrv(srv *http.Server) error {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Err(err)
		}
		close(idleConnsClosed)
	}()

	log.Info().Msg("starting previewer server on" + srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	return nil
}
