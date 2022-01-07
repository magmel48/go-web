package main

import (
	"github.com/magmel48/go-web/internal/app"
	"github.com/magmel48/go-web/internal/config"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	config.Parse()

	log.Printf("starting on %s with %s as base url\n", config.AppDomain, config.BaseShortenerURL)

	shortenerApp := app.NewApp(config.BaseShortenerURL)
	err := fasthttp.ListenAndServe(config.AppDomain, shortenerApp.HTTPHandler())

	if err != nil {
		panic(err)
	}
}
