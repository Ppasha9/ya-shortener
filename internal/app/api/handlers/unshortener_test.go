package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Ppasha9/ya-shortener/internal/app/api"
	"github.com/Ppasha9/ya-shortener/internal/app/auth"
	"github.com/Ppasha9/ya-shortener/internal/app/config"
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnShortenerHandler(t *testing.T) {
	st, err := storage.NewInMemoryStorage(*config.FileStoragePath)
	require.NoError(t, err)

	ctx := context.Background()

	tests := []struct {
		name       string
		reqMethod  string
		reqURLID   string
		userID     uint32
		origURL    string
		respCode   int
		isPositive bool
	}{
		{
			name:       "invalid request method",
			reqMethod:  http.MethodPost,
			reqURLID:   "unknown_url_id",
			userID:     1,
			respCode:   http.StatusMethodNotAllowed,
			isPositive: false,
		},
		{
			name:       "valid request method, unknown url id",
			reqMethod:  http.MethodGet,
			reqURLID:   "unknown_url_id",
			userID:     1,
			respCode:   http.StatusInternalServerError,
			isPositive: false,
		},
		{
			name:       "valid request method, known url id -> 307 redirect",
			reqMethod:  http.MethodGet,
			reqURLID:   "known_url_id",
			userID:     1,
			origURL:    "https://yandex.ru",
			respCode:   http.StatusTemporaryRedirect,
			isPositive: true,
		},
	}

	for _, test := range tests {
		st.Clear()

		t.Run(test.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

			auth, err := auth.NewAuth()
			require.NoError(t, err)

			// инициализируем api
			r := chi.NewRouter()
			api := api.NewAPI(r, st, auth, logger)
			h := NewHandlers(api)
			h.ConfigureRouter()

			if test.origURL != "" {
				st.SaveURL(ctx, test.userID, test.reqURLID, test.origURL)
			}

			request, _ := http.NewRequest(test.reqMethod, "/"+test.reqURLID, nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			authCookie, err := auth.NewJWT(test.userID)
			require.NoError(t, err)
			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    authCookie,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
			api.Router.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			require.Equal(t, test.respCode, res.StatusCode)

			// проверяем значение хэдэра Location в ответе
			if test.isPositive {
				resLoc := res.Header.Get("Location")
				require.NotEmpty(t, resLoc)
				assert.Equal(t, test.origURL, resLoc)
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    "",
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   -1,
			})
		})
	}
}
