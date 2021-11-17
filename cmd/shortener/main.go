package main

import (
	"github.com/magmel48/go-web/internal/app"
	"github.com/valyala/fasthttp"
	"os"
)

func main() {
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "8080"
	}

	shortenerApp := app.NewApp(host, port)
	err := fasthttp.ListenAndServe(host+":"+port, shortenerApp.HandleHTTPRequests)

	if err != nil {
		panic(err)
	}
}
