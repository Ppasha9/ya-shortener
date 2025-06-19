package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	serviceerrors "github.com/Ppasha9/ya-shortener/internal/app/errors"
)

func (h *handlers) UrlsHandler(w http.ResponseWriter, r *http.Request) {
	h.api.Logger.Info("Incoming GET urls request")

	userID, err := h.getUserIDFromCookie(&w, r, h.api.Logger)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrInvalidAuthCookieBytesLen) {
			h.api.Logger.Error("Auth cookie doesn't contain user id")
			http.Error(w, "Auth cookie doesn't contain user id", http.StatusUnauthorized)
			return
		}

		h.api.Logger.Error("Failed to get userID from auth cookie", "err", err.Error())
		http.Error(w, "Failed to get userID from auth cookie", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodGet {
		// Принимаем только GET запросы
		h.api.Logger.Error("Invalid method", "method", r.Method)
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	h.api.Logger.Info("Getting user's urls", "userID", userID)

	userURLs, err := h.api.Service.GetUserURLs(r.Context(), userID)
	statusCode := http.StatusOK
	if err != nil {
		h.api.Logger.Error("Failed to get user's urls", "userID", userID, "err", err.Error())
		http.Error(w, "Failed to get user's urls", http.StatusInternalServerError)
		return
	}

	h.api.Logger.Info(fmt.Sprintf("Got %d urls for user", len(userURLs)), "userID", userID)
	if len(userURLs) == 0 {
		statusCode = http.StatusNoContent
	}

	respBody, err := json.Marshal(&userURLs)
	if err != nil {
		h.api.Logger.Error("Failed to marshal response body", "err", err.Error())
		http.Error(w, "Failed to marshal response body", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(respBody)
}
