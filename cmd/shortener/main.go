package main

import (
	"github.com/magmel48/go-web/internal/app"
	"github.com/valyala/fasthttp"
	"os"
)

func main() {
	address := os.Getenv("SERVER_ADDRESS")

	if address == "" {
		address = "localhost:8080"
	}

	shortenerApp := app.NewApp("http://" + address)
	err := fasthttp.ListenAndServe(address, shortenerApp.HTTPHandler())

	if err != nil {
		panic(err)
	}
}
