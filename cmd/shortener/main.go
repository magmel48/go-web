package main

import (
	"fmt"
	"github.com/magmel48/go-web/internal/app"
	"github.com/magmel48/go-web/internal/config"
	"github.com/magmel48/go-web/internal/db"
	"github.com/valyala/fasthttp"
)

func main() {
	config.Parse()
	db.Connect()

	fmt.Printf("starting on %s with %s as base url\n", config.AppDomain, config.BaseShortenerURL)

	shortenerApp := app.NewApp(config.BaseShortenerURL)
	err := fasthttp.ListenAndServe(config.AppDomain, shortenerApp.HTTPHandler())

	if err != nil {
		panic(err)
	}
}
