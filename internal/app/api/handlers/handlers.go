package handlers

import (
	"net/http"

	"github.com/Ppasha9/ya-shortener/internal/app/api"
	"github.com/Ppasha9/ya-shortener/internal/app/api/middleware"
)

type handlers struct {
	api *api.API
}

func NewHandlers(a *api.API) *handlers {
	return &handlers{
		api: a,
	}
}

func (h *handlers) ConfigureRouter() {
	// Хэндлеры первого спринта
	h.api.Router.Handle("/", middleware.WithLogging(middleware.WithCompress(http.HandlerFunc(h.ShortenerHandler)), h.api.Logger))
	h.api.Router.Handle("/{id}", middleware.WithLogging(middleware.WithCompress(http.HandlerFunc(h.UnShortenerHandler)), h.api.Logger))

	// Хэндлеры второго спринта
	h.api.Router.Handle("/api/shorten", middleware.WithLogging(middleware.WithCompress(http.HandlerFunc(h.ShortenHandler)), h.api.Logger))
}
