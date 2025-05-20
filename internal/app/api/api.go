package api

import (
	"log/slog"

	"github.com/go-chi/chi"

	"github.com/Ppasha9/ya-shortener/internal/app/service"
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
)

type API struct {
	Router  *chi.Mux
	Service *service.Service
	Logger  *slog.Logger
}

func NewAPI(r *chi.Mux, s *storage.InMemoryStorage, l *slog.Logger) *API {
	api := &API{
		Router:  r,
		Service: service.NewService(s),
		Logger:  l,
	}
	return api
}
