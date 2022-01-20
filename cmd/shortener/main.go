package main

import (
	"context"
	"github.com/magmel48/go-web/internal/app"
	"github.com/magmel48/go-web/internal/config"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		termChan := make(chan os.Signal, 1)
		signal.Notify(termChan, os.Interrupt, syscall.SIGTERM)
		cancel()
	}()

	log.Printf("starting on %s with %s as base url\n", config.AppDomain, config.BaseShortenerURL)

	shortenerApp := app.NewApp(ctx, config.BaseShortenerURL)
	err := fasthttp.ListenAndServe(config.AppDomain, shortenerApp.HTTPHandler())

	if err != nil {
		panic(err)
	}
}
