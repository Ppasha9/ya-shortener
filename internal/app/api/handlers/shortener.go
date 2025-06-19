package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Ppasha9/ya-shortener/internal/app/config"
	servicecontext "github.com/Ppasha9/ya-shortener/internal/app/context"
	serviceerrors "github.com/Ppasha9/ya-shortener/internal/app/errors"
	"github.com/Ppasha9/ya-shortener/internal/app/model"
)

func isValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func (h *handlers) ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	h.api.Logger.Info("Incoming POST shortener request")

	userID, err := servicecontext.GetUserID(r.Context())
	if err != nil {
		h.api.Logger.Error("Failed to get userID from auth cookie", "err", err.Error())
		http.Error(w, "Failed to get userID from auth cookie", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		// Принимаем только POST запросы
		h.api.Logger.Error("Invalid method", "method", r.Method)
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем наличие хэдэра Content-Type и его значение
	ctHeader := r.Header.Get("Content-Type")
	if !strings.Contains(ctHeader, "text/plain") && !strings.Contains(ctHeader, "application/x-gzip") {
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

	h.api.Logger.Info("Generating short url", "userID", userID, "orig_url", origURL)
	statusCode := http.StatusCreated
	shortURL, err := h.api.Service.MakeShortURL(r.Context(), userID, origURL)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrOrigURLDuplicate) {
			h.api.Logger.Error("Trying to shorten certain original url one more time", "userID", userID, "req_body", origURL)
			statusCode = http.StatusConflict
		} else {
			h.api.Logger.Error("Failed to generate short url", "userID", userID, "req_body", origURL, "err", err.Error())
			http.Error(w, "Failed to generate short url", http.StatusInternalServerError)
			return
		}
	}

	shortURL = *config.BaseURL + "/" + shortURL

	h.api.Logger.Info("Generated short url", "userID", userID, "orig_url", origURL, "short_url", shortURL)

	w.WriteHeader(statusCode)
	w.Write([]byte(shortURL))
}

func (h *handlers) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	h.api.Logger.Info("Incoming POST shorten request")

	userID, err := servicecontext.GetUserID(r.Context())
	if err != nil {
		h.api.Logger.Error("Failed to get userID from auth cookie", "err", err.Error())
		http.Error(w, "Failed to get userID from auth cookie", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		// Принимаем только POST запросы
		h.api.Logger.Error("Invalid method", "method", r.Method)
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем наличие хэдэра Content-Type и его значение
	ctHeader := r.Header.Get("Content-Type")
	if !strings.Contains(ctHeader, "application/json") && !strings.Contains(ctHeader, "application/x-gzip") {
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
		h.api.Logger.Error("URL from request isn't valid", "userID", userID, "req_url", origURL)
		http.Error(w, "URL from request isn't valid", http.StatusBadRequest)
		return
	}

	h.api.Logger.Info("Generating short url", "userID", userID, "orig_url", origURL)
	statusCode := http.StatusCreated
	shortURL, err := h.api.Service.MakeShortURL(r.Context(), userID, origURL)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrOrigURLDuplicate) {
			h.api.Logger.Error("Trying to shorten certain original url one more time", "userID", userID, "req_body", origURL)
			statusCode = http.StatusConflict
		} else {
			h.api.Logger.Error("Failed to generate short url", "userID", userID, "req_url", origURL, "err", err.Error())
			http.Error(w, "Failed to generate short url", http.StatusInternalServerError)
			return
		}
	}

	shortURL = *config.BaseURL + "/" + shortURL

	h.api.Logger.Info("Generated short url", "userID", userID, "orig_url", origURL, "short_url", shortURL)

	resp := model.ShortenResponse{
		Result: shortURL,
	}
	respBody, err := json.Marshal(&resp)
	if err != nil {
		h.api.Logger.Error("Failed to marshal response body", "err", err.Error())
		http.Error(w, "Failed to marshal response body", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(respBody)
}

func (h *handlers) ShortenBatchHandler(w http.ResponseWriter, r *http.Request) {
	h.api.Logger.Info("Incoming POST shorten batch request")

	userID, err := servicecontext.GetUserID(r.Context())
	if err != nil {
		h.api.Logger.Error("Failed to get userID from auth cookie", "err", err.Error())
		http.Error(w, "Failed to get userID from auth cookie", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		// Принимаем только POST запросы
		h.api.Logger.Error("Invalid method", "method", r.Method)
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем наличие хэдэра Content-Type и его значение
	ctHeader := r.Header.Get("Content-Type")
	if !strings.Contains(ctHeader, "application/json") && !strings.Contains(ctHeader, "application/x-gzip") {
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

	var req []model.ShortenBatchItemRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		h.api.Logger.Error("Failed to unmarshal request body", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range req {
		if !isValidURL(v.OrigURL) {
			h.api.Logger.Error("URL from request isn't valid", "userID", userID, "original_url", v.OrigURL, "correlation_id", v.CorrID)
			http.Error(w, "URL from request isn't valid", http.StatusBadRequest)
			return
		}
	}

	h.api.Logger.Info(fmt.Sprintf("Generating %d short urls", len(req)), "userID", userID)
	resp, err := h.api.Service.MakeShortURLsBatch(r.Context(), userID, req)
	if err != nil {
		h.api.Logger.Error("Failed to generate batch of short urls", "userID", userID, "batch_size", len(req), "err", err.Error())
		http.Error(w, "Failed to generate batch of short urls", http.StatusInternalServerError)
		return
	}

	for i := range resp {
		resp[i].ShortURL = *config.BaseURL + "/" + resp[i].ShortURL
	}

	h.api.Logger.Info("Generated batch of short urls", "userID", userID, "batch_size", len(req))

	respBody, err := json.Marshal(&resp)
	if err != nil {
		h.api.Logger.Error("Failed to marshal response body", "err", err.Error())
		http.Error(w, "Failed to marshal response body", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}
