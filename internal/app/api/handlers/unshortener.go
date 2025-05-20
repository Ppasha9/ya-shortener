package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (h *handlers) UnShortenerHandler(w http.ResponseWriter, r *http.Request) {
	h.api.Logger.Info("Incoming GET unshortener request")

	if r.Method != http.MethodGet {
		// Принимаем только GET запросы
		h.api.Logger.Error("Invalid method", "method", r.Method)
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	urlID := chi.URLParam(r, "id")
	if urlID == "" {
		h.api.Logger.Error("Invalid url id", "url_id", urlID)
		http.Error(w, "Invalid url id", http.StatusBadRequest)
		return
	}

	h.api.Logger.Info("Getting original url by url_id", "url_id", urlID)

	origURL, err := h.api.Service.GetOriginalURL(urlID)
	if err != nil {
		h.api.Logger.Error("Failed to get original url by url id", "err", err.Error())
		http.Error(w, "Failed to get original url by url id", http.StatusInternalServerError)
		return
	}

	h.api.Logger.Info("Got original url by url_id", "url_id", urlID, "orig_url", origURL)

	w.Header().Add("Location", origURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
