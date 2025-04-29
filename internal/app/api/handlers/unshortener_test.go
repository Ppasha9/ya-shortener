package handlers

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Ppasha9/ya-shortener/internal/app/api"
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnShortenerHandler(t *testing.T) {
	db := storage.NewDatabase()

	tests := []struct {
		name       string
		reqMethod  string
		reqUrlID   string
		origUrl    string
		respCode   int
		isPositive bool
	}{
		{
			name:       "invalid request method",
			reqMethod:  http.MethodPost,
			reqUrlID:   "unknown_url_id",
			respCode:   http.StatusBadRequest,
			isPositive: false,
		},
		{
			name:       "valid request method, unknown url id",
			reqMethod:  http.MethodGet,
			reqUrlID:   "unknown_url_id",
			respCode:   http.StatusBadRequest,
			isPositive: false,
		},
		{
			name:       "valid request method, known url id -> 307 redirect",
			reqMethod:  http.MethodGet,
			reqUrlID:   "known_url_id",
			origUrl:    "https://yandex.ru",
			respCode:   http.StatusTemporaryRedirect,
			isPositive: true,
		},
	}

	for _, test := range tests {
		db.Clear()

		t.Run(test.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

			// инициализируем api
			r := chi.NewRouter()
			api := api.NewApi(r, db, logger)
			h := NewHandlers(api)
			h.ConfigureRouter()

			if test.origUrl != "" {
				db.SaveURL(test.reqUrlID, test.origUrl)
			}

			request, _ := http.NewRequest(test.reqMethod, "/"+test.reqUrlID, nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			api.Router.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			require.Equal(t, test.respCode, res.StatusCode)

			// проверяем значение хэдэра Location в ответе
			if test.isPositive {
				resLoc := res.Header.Get("Location")
				require.NotEmpty(t, resLoc)
				assert.Equal(t, test.origUrl, resLoc)
			}
		})
	}
}
