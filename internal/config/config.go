package config

import (
	"flag"
	"os"
	"strings"
)

var defaultProtocol = "http://"

var (
	// Address where the server starts their job
	Address string
	// AppDomain is parsed from Address - where the server starts their job
	AppDomain string
	// BaseShortenerURL is what will be put in short link as domain ({domain}/{short_link_id})
	BaseShortenerURL string
	// FilePath is deprecated, previously was using for storing backup
	FilePath string
	// SecretKey is secret character sequence that is using for encoding/decoding user identifiers
	SecretKey string
	// DatabaseDSN is database connection string
	DatabaseDSN string
)

// Parse parses flags and gets default values for them from environment variables. Hides details of envs ingestion.
func Parse() {
	flag.StringVar(&Address, "a", os.Getenv("SERVER_ADDRESS"), "server address")
	flag.StringVar(&BaseShortenerURL, "b", os.Getenv("BASE_URL"), "base url for shortened urls")
	flag.StringVar(&FilePath, "f", os.Getenv("FILE_STORAGE_PATH"), "file path for shortened links")
	flag.StringVar(&SecretKey, "s", os.Getenv("SECRET_KEY"), "secret key for sessions")
	flag.StringVar(&DatabaseDSN, "d", os.Getenv("DATABASE_DSN"), "database connection string")
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
		FilePath = "backup.txt"
	}

	if SecretKey == "" {
		SecretKey = "secret_key"
	}
}
