package handlers

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Ppasha9/ya-shortener/internal/app/api"
	"github.com/Ppasha9/ya-shortener/internal/app/config"
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnShortenerHandler(t *testing.T) {
	st, err := storage.NewInMemoryStorage(*config.FileStoragePath)
	require.NoError(t, err)

	tests := []struct {
		name       string
		reqMethod  string
		reqURLID   string
		origURL    string
		respCode   int
		isPositive bool
	}{
		{
			name:       "invalid request method",
			reqMethod:  http.MethodPost,
			reqURLID:   "unknown_url_id",
			respCode:   http.StatusMethodNotAllowed,
			isPositive: false,
		},
		{
			name:       "valid request method, unknown url id",
			reqMethod:  http.MethodGet,
			reqURLID:   "unknown_url_id",
			respCode:   http.StatusInternalServerError,
			isPositive: false,
		},
		{
			name:       "valid request method, known url id -> 307 redirect",
			reqMethod:  http.MethodGet,
			reqURLID:   "known_url_id",
			origURL:    "https://yandex.ru",
			respCode:   http.StatusTemporaryRedirect,
			isPositive: true,
		},
	}

	for _, test := range tests {
		st.Clear()

		t.Run(test.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

			// инициализируем api
			r := chi.NewRouter()
			api := api.NewAPI(r, st, logger)
			h := NewHandlers(api)
			h.ConfigureRouter()

			if test.origURL != "" {
				st.SaveURL(test.reqURLID, test.origURL)
			}

			request, _ := http.NewRequest(test.reqMethod, "/"+test.reqURLID, nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()
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
		})
	}
}
