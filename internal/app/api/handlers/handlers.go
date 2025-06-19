package handlers

import (
	"net/http"

	"github.com/Ppasha9/ya-shortener/internal/app/api"
	m "github.com/Ppasha9/ya-shortener/internal/app/api/middleware"
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
	h.api.Router.Handle("/", m.WithAuthCheck(m.WithLogging(m.WithCompress(http.HandlerFunc(h.ShortenerHandler)), h.api.Logger), h.api.Auth, h.api.Logger))
	h.api.Router.Handle("/{id}", m.WithLogging(m.WithCompress(http.HandlerFunc(h.UnShortenerHandler)), h.api.Logger))

	// Хэндлеры второго спринта
	h.api.Router.Handle("/api/shorten", m.WithAuthCheck(m.WithLogging(m.WithCompress(http.HandlerFunc(h.ShortenHandler)), h.api.Logger), h.api.Auth, h.api.Logger))

	// Хэндлеры третьего спринта
	h.api.Router.Handle("/ping", m.WithLogging(http.HandlerFunc(h.PingHandler), h.api.Logger))
	h.api.Router.Handle("/api/shorten/batch", m.WithAuthCheck(m.WithLogging(m.WithCompress(http.HandlerFunc(h.ShortenBatchHandler)), h.api.Logger), h.api.Auth, h.api.Logger))

	// Хэндлеры четвертого спринта
	h.api.Router.Handle("/api/user/urls", m.WithAuthCheck(m.WithLogging(m.WithCompress(http.HandlerFunc(h.UrlsHandler)), h.api.Logger), h.api.Auth, h.api.Logger))
}
