package api

import (
	"log/slog"

	"github.com/go-chi/chi"

	"github.com/Ppasha9/ya-shortener/internal/app/auth"
	"github.com/Ppasha9/ya-shortener/internal/app/service"
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
)

type API struct {
	Router  *chi.Mux
	Service *service.Service
	Auth    *auth.Auth
	Logger  *slog.Logger
}

func NewAPI(r *chi.Mux, s storage.Storage, a *auth.Auth, l *slog.Logger) *API {
	api := &API{
		Router:  r,
		Service: service.NewService(s),
		Auth:    a,
		Logger:  l,
	}
	return api
}
