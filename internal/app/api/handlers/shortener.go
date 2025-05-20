package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/Ppasha9/ya-shortener/internal/app/config"
)

func (h *handlers) ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	h.api.Logger.Info("Incoming POST shortener request")

	if r.Method != http.MethodPost {
		// Принимаем только POST запросы
		h.api.Logger.Error("Invalid method", "method", r.Method)
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем наличие хэдэра Content-Type и его значение
	ctHeader := r.Header.Get("Content-Type")
	if !strings.Contains(ctHeader, "text/plain") {
		h.api.Logger.Error("Invalid Content-Type header", "header_value", ctHeader)
		http.Error(w, "Invalid Content-Type header", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.api.Logger.Error("Failed to read request body", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	origURL := string(body)
	if !strings.HasPrefix(origURL, "http://") && !strings.HasPrefix(origURL, "https://") {
		h.api.Logger.Error("Request body is not an url", "req_body", origURL)
		http.Error(w, "Request body is not an url", http.StatusBadRequest)
		return
	}

	h.api.Logger.Info("Generating short url", "orig_url", origURL)
	shortURL, err := h.api.Service.MakeShortURL(origURL)
	if err != nil {
		h.api.Logger.Error("Failed to generate short url", "req_body", origURL)
		http.Error(w, "Failed to generate short url", http.StatusInternalServerError)
		return
	}

	shortURL = *config.BaseURL + "/" + shortURL

	h.api.Logger.Info("Generated short url", "orig_url", origURL, "short_url", shortURL)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}
