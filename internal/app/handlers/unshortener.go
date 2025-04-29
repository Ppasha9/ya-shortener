package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Ppasha9/ya-shortener/internal/app/storage"
)

type UnShortenerHandler struct {
	Logger *slog.Logger
}

func (h UnShortenerHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Incoming GET unshortener request")

	if r.Method != http.MethodGet {
		// Принимаем только GET запросы
		h.Logger.Error("Invalid method", "method", r.Method)
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	queryParams := mux.Vars(r)
	urlID := queryParams["id"]
	if urlID == "" {
		h.Logger.Error("Invalid url id", "url_id", urlID)
		http.Error(w, "Invalid url id", http.StatusBadRequest)
		return
	}

	h.Logger.Info("Getting original url by url_id", "url_id", urlID)

	origURL, err := storage.GetOriginalURL(urlID)
	if err != nil {
		h.Logger.Error("Failed to get original url by url id", "err", err.Error())
		http.Error(w, "Failed to get original url by url id", http.StatusBadRequest)
		return
	}

	h.Logger.Info("Got original url by url_id", "url_id", urlID, "orig_url", origURL)

	w.Header().Add("Location", origURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func NewUnShortenerHandler(logger *slog.Logger) UnShortenerHandler {
	return UnShortenerHandler{
		Logger: logger,
	}
}
