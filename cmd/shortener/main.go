package main

import (
	"github.com/magmel48/go-web/internal/app"
	"net/http"
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
	http.HandleFunc("/", shortenerApp.HandleHTTPRequests)

	err := http.ListenAndServe(host+":"+port, nil)
	if err != nil {
		panic(err)
	}
}
