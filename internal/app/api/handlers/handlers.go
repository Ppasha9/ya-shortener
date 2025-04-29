package handlers

import (
	"github.com/Ppasha9/ya-shortener/internal/app/api"
)

type handlers struct {
	api *api.Api
}

func NewHandlers(a *api.Api) *handlers {
	return &handlers{
		api: a,
	}
}

func (h *handlers) ConfigureRouter() {
	h.api.Router.HandleFunc("/", h.ShortenerHandler)
	h.api.Router.HandleFunc("/{id}", h.UnShortenerHandler)
}
