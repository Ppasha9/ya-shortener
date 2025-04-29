package api

import (
	"log/slog"

	"github.com/go-chi/chi"

	"github.com/Ppasha9/ya-shortener/internal/app/storage"
)

type Api struct {
	Router  *chi.Mux
	Storage *storage.Database
	Logger  *slog.Logger
}

func NewApi(r *chi.Mux, s *storage.Database, l *slog.Logger) *Api {
	api := &Api{
		Router:  r,
		Storage: s,
		Logger:  l,
	}
	return api
}
