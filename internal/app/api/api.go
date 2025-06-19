package api

import (
	"log/slog"

	"github.com/go-chi/chi"

	"github.com/Ppasha9/ya-shortener/internal/app/crypt"
	"github.com/Ppasha9/ya-shortener/internal/app/service"
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
)

type API struct {
	Router  *chi.Mux
	Service *service.Service
	Crypt   *crypt.Crypt
	Logger  *slog.Logger
}

func NewAPI(r *chi.Mux, s storage.Storage, c *crypt.Crypt, l *slog.Logger) *API {
	api := &API{
		Router:  r,
		Service: service.NewService(s),
		Crypt:   c,
		Logger:  l,
	}
	return api
}
