package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Ppasha9/ya-shortener/internal/app/auth"
	servicecontext "github.com/Ppasha9/ya-shortener/internal/app/context"
)

var authCookieName = "access_token"

func WithAuthCheck(h http.Handler, authService *auth.Auth, logger *slog.Logger) http.Handler {
	authCheckFn := func(w http.ResponseWriter, r *http.Request) {
		var userID uint32
		var err error
		shouldCreateJWT := false

		authCookie, err := r.Cookie(authCookieName)
		if err != nil {
			shouldCreateJWT = true
		} else {
			userID, err = authService.ParseUserIDFromJWT(authCookie.Value)
			if err != nil {
				shouldCreateJWT = true
			}
		}

		if shouldCreateJWT {
			userID, err = auth.GenerateUserID()
			if err != nil {
				logger.Error("Failed to generate user id", "err", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			token, err := authService.NewJWT(userID)
			if err != nil {
				logger.Error("Failed to generate jwt", "err", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     authCookieName,
				Value:    token,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Expires:  time.Now().Add(authService.TokenTTL),
			})
		}

		ctx := servicecontext.WithUserID(r.Context(), userID)
		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(authCheckFn)
}
