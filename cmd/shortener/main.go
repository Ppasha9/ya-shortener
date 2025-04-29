package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Ppasha9/ya-shortener/internal/app/handlers"
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
	"github.com/gorilla/mux"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	storage.Init()

	shortenerHandler := handlers.NewShortenerHandler(logger)
	unshortenerHandler := handlers.NewUnShortenerHandler(logger)

	r := mux.NewRouter()
	r.HandleFunc("/", shortenerHandler.ServerHTTP)
	r.HandleFunc("/{id}", unshortenerHandler.ServerHTTP)

	logger.Info("Starting shortener...")
	err := http.ListenAndServe(`:8080`, r)
	logger.Info("Stopping shortener...")
	return err
}
