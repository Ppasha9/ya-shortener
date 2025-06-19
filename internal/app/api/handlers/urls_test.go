package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Ppasha9/ya-shortener/internal/app/api"
	"github.com/Ppasha9/ya-shortener/internal/app/auth"
	"github.com/Ppasha9/ya-shortener/internal/app/config"
	"github.com/Ppasha9/ya-shortener/internal/app/model"
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
)

func TestUrlsHandler(t *testing.T) {
	st, err := storage.NewInMemoryStorage(*config.FileStoragePath)
	require.NoError(t, err)

	ctx := context.Background()

	tests := []struct {
		name       string
		reqMethod  string
		userID     uint32
		userURLs   []model.URLsPair
		respCode   int
		isPositive bool
	}{
		{
			name:       "invalid request method",
			reqMethod:  http.MethodPost,
			userID:     1,
			respCode:   http.StatusMethodNotAllowed,
			isPositive: false,
		},
		{
			name:       "valid request method, user has no urls",
			reqMethod:  http.MethodGet,
			userID:     1,
			respCode:   http.StatusNoContent,
			isPositive: false,
		},
		{
			name:      "valid request method, known url id -> 307 redirect",
			reqMethod: http.MethodGet,
			userID:    1,
			userURLs: []model.URLsPair{
				{
					ShortURL: "dfuy45bn",
					OrigURL:  "https://yandex.ru",
				},
				{
					ShortURL: "ls53ot65",
					OrigURL:  "https://practicum.yandex.ru",
				},
			},
			respCode:   http.StatusOK,
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

			if len(test.userURLs) > 0 {
				st.SaveURLs(ctx, test.userID, test.userURLs)
			}

			request, _ := http.NewRequest(test.reqMethod, "/api/user/urls", nil)

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
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)

				var resp []model.URLsPair
				err = json.Unmarshal(resBody, &resp)
				require.NoError(t, err)

				require.Equal(t, len(resp), len(test.userURLs))
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
