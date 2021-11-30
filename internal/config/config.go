package config

import (
	"flag"
	"os"
	"strings"
)

var defaultProtocol = "http://"

var (
	Address          string
	AppDomain        string
	BaseShortenerURL string
	FilePath         string
)

// Parse parses flags and gets default values for them from environment variables. Hides details of envs ingestion.
func Parse() {
	flag.StringVar(&Address, "a", os.Getenv("SERVER_ADDRESS"), "server address")
	flag.StringVar(&BaseShortenerURL, "b", os.Getenv("BASE_URL"), "base url for shortened urls")
	flag.StringVar(&FilePath, "f", os.Getenv("FILE_STORAGE_PATH"), "file path for shortened links")
	flag.Parse()

	if Address == "" {
		Address = "localhost:8080"
	}

	addressComponents := strings.Split(Address, defaultProtocol)
	if len(addressComponents) > 1 {
		AppDomain = addressComponents[1]
	} else {
		AppDomain = addressComponents[0]
	}

	if BaseShortenerURL == "" {
		BaseShortenerURL = Address

		if strings.Index(BaseShortenerURL, defaultProtocol) != 0 {
			BaseShortenerURL = defaultProtocol + BaseShortenerURL
		}
	}

	if FilePath == "" {
		FilePath = "links.txt"
	}
}
