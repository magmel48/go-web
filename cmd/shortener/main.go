package main

import (
	"context"
	"github.com/magmel48/go-web/internal/app"
	"github.com/magmel48/go-web/internal/config"
	"github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	shortenerApp := app.NewApp(ctx, config.BaseShortenerURL)
	server := fasthttp.Server{
		Handler:         shortenerApp.HTTPHandler(),
		CloseOnShutdown: true,
		IdleTimeout:     time.Second,
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		log.Printf("starting on %s with %s as base url\n", config.AppDomain, config.BaseShortenerURL)
		return server.ListenAndServe(config.AppDomain)
	})

	eg.Go(func() error {
		<-ctx.Done()

		log.Println("stopping the service...")

		err := server.Shutdown()
		return err
	})

	log.Println(eg.Wait())
}
