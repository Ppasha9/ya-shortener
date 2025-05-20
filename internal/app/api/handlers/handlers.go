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
	h.api.Router.Handle("/", middleware.WithLogging(http.HandlerFunc(h.ShortenerHandler), h.api.Logger))
	h.api.Router.Handle("/{id}", middleware.WithLogging(http.HandlerFunc(h.UnShortenerHandler), h.api.Logger))
}
