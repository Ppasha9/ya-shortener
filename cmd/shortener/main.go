package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"github.com/Ppasha9/ya-shortener/internal/app/api"
	"github.com/Ppasha9/ya-shortener/internal/app/api/handlers"
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	s := storage.NewDatabase()

	r := chi.NewRouter()
	api := api.NewAPI(r, s, logger)
	h := handlers.NewHandlers(api)
	h.ConfigureRouter()

	logger.Info("Starting shortener...")
	err := http.ListenAndServe(`:8080`, r)
	logger.Info("Stopping shortener...")
	return err
}
