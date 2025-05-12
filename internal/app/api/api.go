package api

import (
	"log/slog"

	"github.com/go-chi/chi"

	"github.com/Ppasha9/ya-shortener/internal/app/storage"
)

type API struct {
	Router  *chi.Mux
	Storage *storage.Database
	Logger  *slog.Logger
}

func NewAPI(r *chi.Mux, s *storage.Database, l *slog.Logger) *API {
	api := &API{
		Router:  r,
		Storage: s,
		Logger:  l,
	}
	return api
}
