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
	"github.com/Ppasha9/ya-shortener/internal/app/auth"
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

	auth, err := auth.NewAuth()
	if err != nil {
		logger.Error("cannot initialize auth crypt", "err", err.Error())
		return err
	}

	var a *api.API
	r := chi.NewRouter()

	if *config.DatabaseDSN == "" {
		memStorage, err := storage.NewInMemoryStorage(*config.FileStoragePath)
		if err != nil {
			logger.Error("cannot open file storage file", "err", err.Error())
			return err
		}
		a = api.NewAPI(r, memStorage, auth, logger)
	} else {
		dbStorage, err := storage.NewDatabase(*config.DatabaseDSN)
		if err != nil {
			logger.Error("cannot open db storage connection", "err", err.Error())
			return err
		}
		a = api.NewAPI(r, dbStorage, auth, logger)
	}
	defer a.Service.Storage.Close()

	h := handlers.NewHandlers(a)
	h.ConfigureRouter()

	logger.Info(fmt.Sprintf("Starting shortener on %s ...", *config.ServerAddr))
	err = http.ListenAndServe(*config.ServerAddr, r)
	logger.Info("Stopping shortener...")

	return err
}
