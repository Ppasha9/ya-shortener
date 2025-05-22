package config

import (
	"flag"
	"os"
)

var (
	ServerAddr      = flag.String("a", ":8080", "address and port to run server")
	BaseURL         = flag.String("b", "http://localhost:8080", "base url for returning short urls")
	FileStoragePath = flag.String("f", "storage.json", "path to the urls storage file")
)

func ParseArgs() {
	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		ServerAddr = &envServerAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		BaseURL = &envBaseURL
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FileStoragePath = &envFileStoragePath
	}
}
