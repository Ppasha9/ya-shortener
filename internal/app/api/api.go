package api

import (
	"log/slog"

	"github.com/go-chi/chi"

	"github.com/Ppasha9/ya-shortener/internal/app/config"
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
)

type API struct {
	Router  *chi.Mux
	Storage *storage.Database
	Config  config.Config
	Logger  *slog.Logger
}

func NewAPI(r *chi.Mux, s *storage.Database, c config.Config, l *slog.Logger) *API {
	api := &API{
		Router:  r,
		Storage: s,
		Config:  c,
		Logger:  l,
	}
	return api
}
