package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"github.com/Ppasha9/ya-shortener/internal/app/api"
	"github.com/Ppasha9/ya-shortener/internal/app/api/handlers"
	"github.com/Ppasha9/ya-shortener/internal/app/config"
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flag.Parse()
	config.ParseArgs()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	s := storage.NewInMemoryStorage()

	r := chi.NewRouter()
	api := api.NewAPI(r, s, logger)
	h := handlers.NewHandlers(api)
	h.ConfigureRouter()

	logger.Info(fmt.Sprintf("Starting shortener on %s ...", *config.ServerAddr))
	err := http.ListenAndServe(*config.ServerAddr, r)
	logger.Info("Stopping shortener...")
	return err
}
