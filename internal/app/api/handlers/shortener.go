package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Ppasha9/ya-shortener/internal/app/config"
	"github.com/Ppasha9/ya-shortener/internal/app/model"
)

func isValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

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
		http.Error(w, "Invalid Content-Type header", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.api.Logger.Error("Failed to read request body", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	origURL := string(body)
	if !isValidURL(origURL) {
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

func (h *handlers) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	h.api.Logger.Info("Incoming POST shorten request")

	if r.Method != http.MethodPost {
		// Принимаем только POST запросы
		h.api.Logger.Error("Invalid method", "method", r.Method)
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем наличие хэдэра Content-Type и его значение
	ctHeader := r.Header.Get("Content-Type")
	if !strings.Contains(ctHeader, "application/json") {
		h.api.Logger.Error("Invalid Content-Type header", "header_value", ctHeader)
		http.Error(w, "Invalid Content-Type header", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.api.Logger.Error("Failed to read request body", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req model.ShortenRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		h.api.Logger.Error("Failed to unmarshal request body", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	origURL := req.URL
	if !isValidURL(origURL) {
		h.api.Logger.Error("URL from request isn't valid", "req_url", origURL)
		http.Error(w, "URL from request isn't valid", http.StatusBadRequest)
		return
	}

	h.api.Logger.Info("Generating short url", "orig_url", origURL)
	shortURL, err := h.api.Service.MakeShortURL(origURL)
	if err != nil {
		h.api.Logger.Error("Failed to generate short url", "req_url", origURL)
		http.Error(w, "Failed to generate short url", http.StatusInternalServerError)
		return
	}

	shortURL = *config.BaseURL + "/" + shortURL

	h.api.Logger.Info("Generated short url", "orig_url", origURL, "short_url", shortURL)

	resp := model.ShortenResponse{
		Result: shortURL,
	}
	respBody, err := json.Marshal(&resp)
	if err != nil {
		h.api.Logger.Error("Failed to marshal response body", "err", err.Error())
		http.Error(w, "Failed to marshal response body", http.StatusInternalServerError)
		return
	}

	w.Header().Add(`Content-Type`, `application/json`)
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}
