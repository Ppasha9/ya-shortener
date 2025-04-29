package config

import "flag"

type Config struct {
	ServerAddr string
	BaseURL    string
}

func (c Config) Parse() {
	flag.StringVar(&c.ServerAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080/", "base url for returning short urls")
	flag.Parse()
}
