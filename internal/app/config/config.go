package config

import "flag"

var (
	ServerAddr = flag.String("a", ":8080", "address and port to run server")
	BaseURL    = flag.String("b", "http://localhost:8080/", "base url for returning short urls")
)
