package main

import (
	"flag"
	"fmt"
	"github.com/magmel48/go-web/internal/app"
	"github.com/valyala/fasthttp"
	"os"
	"strings"
)

func main() {
	var address string
	var appDomain string
	var baseShortenerURL string

	flag.StringVar(&address,"a", os.Getenv("SERVER_ADDRESS"), "server address")
	flag.StringVar(&baseShortenerURL, "b", os.Getenv("BASE_URL"), "base url for shortened urls")
	flag.Parse()

	if address == "" {
		address = "localhost:8080"
	}

	addressComponents := strings.Split(address, "//")
	if len(addressComponents) > 1 {
		appDomain = addressComponents[1]
	} else {
		appDomain = addressComponents[0]
	}

	if baseShortenerURL == "" {
		baseShortenerURL = address
	}

	fmt.Printf("starting on %s with %s as base url\n", appDomain, baseShortenerURL)

	shortenerApp := app.NewApp(baseShortenerURL)
	err := fasthttp.ListenAndServe(appDomain, shortenerApp.HTTPHandler())

	if err != nil {
		panic(err)
	}
}
