package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Ppasha9/ya-shortener/internal/app/storage"
	"github.com/Ppasha9/ya-shortener/internal/app/urlshortener"
)

type ShortenerHandler struct {
	Logger *slog.Logger
}

func (h ShortenerHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Incoming POST shortener request")

	if r.Method != http.MethodPost {
		// Принимаем только POST запросы
		h.Logger.Error("Invalid method", "method", r.Method)
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	// Проверяем наличие хэдэра Content-Type и его значение
	ctHeader := r.Header.Get("Content-Type")
	if !strings.Contains(ctHeader, "text/plain") {
		h.Logger.Error("Invalid Content-Type header", "header_value", ctHeader)
		http.Error(w, "Invalid Content-Type header", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Error("Failed to read request body", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	origURL := string(body)
	if !strings.HasPrefix(origURL, "http://") && !strings.HasPrefix(origURL, "https://") {
		h.Logger.Error("Request body is not an url", "req_body", origURL)
		http.Error(w, "Request body is not an url", http.StatusBadRequest)
		return
	}

	h.Logger.Info("Generating short url", "orig_url", origURL)

	shortURL := urlshortener.MakeShortURL(origURL)
	storage.SaveURL(shortURL, origURL)

	shortURL = "http://localhost:8080/" + shortURL

	h.Logger.Info("Generated short url", "orig_url", origURL, "short_url", shortURL)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func NewShortenerHandler(logger *slog.Logger) ShortenerHandler {
	return ShortenerHandler{
		Logger: logger,
	}
}
