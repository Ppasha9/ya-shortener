package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Ppasha9/ya-shortener/internal/app/crypt"
	serviceerrors "github.com/Ppasha9/ya-shortener/internal/app/errors"
)

var authCookieName = "auth_cookie"

func (h *handlers) genUserIDAndSetCookie(w *http.ResponseWriter) (uint32, error) {
	userID, err := crypt.GenerateUserID()
	if err != nil {
		return userID, err
	}

	cookieVal := h.api.Crypt.GenerateAuthCookie(userID)
	http.SetCookie(*w, &http.Cookie{
		Name:  authCookieName,
		Value: cookieVal,
		Path:  "/",
	})

	return userID, nil
}

func (h *handlers) getUserIDFromCookie(w *http.ResponseWriter, r *http.Request, logger *slog.Logger) (uint32, error) {
	authCookie, err := r.Cookie(authCookieName)

	// если куки нет, то мы должны сгенерировать новый userID и новую куку
	if err != nil {
		logger.Info("No auth cookie found, generate new.", "err", err.Error())
		userID, err := h.genUserIDAndSetCookie(w)
		if err != nil {
			return userID, err
		}
		return userID, nil
	}

	logger.Info("Auth cookie found.", "val", authCookie.Value)

	// если кука есть, надо проверить ее валидность
	// 1. Если кука невалидна, то нужно сгенерить новый userID и новую куку
	// 2. Если валидна, то достаем из нее userID
	userID, err := h.api.Crypt.GetUserIDFromAuthCookie(authCookie.Value)
	if err != nil && errors.Is(err, serviceerrors.ErrInvalidAuthCookie) {
		userID, err := h.genUserIDAndSetCookie(w)
		if err != nil {
			return userID, err
		}
		return userID, nil
	}
	if err != nil {
		return userID, err
	}

	return userID, nil
}
