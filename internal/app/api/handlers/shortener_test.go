package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Ppasha9/ya-shortener/internal/app/api"
	"github.com/Ppasha9/ya-shortener/internal/app/config"
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
)

func TestShortenerHandler(t *testing.T) {
	*config.BaseURL = "http://baseurl"
	st := storage.NewInMemoryStorage()

	tests := []struct {
		name           string
		reqMethod      string
		reqContentType string
		reqURL         string
		respCode       int
		isPositive     bool
	}{
		{
			name:       "invalid request method",
			reqMethod:  http.MethodGet,
			respCode:   http.StatusMethodNotAllowed,
			isPositive: false,
		},
		{
			name:           "valid request method, but invalid content_type",
			reqMethod:      http.MethodPost,
			reqContentType: "application/json",
			respCode:       http.StatusBadRequest,
			isPositive:     false,
		},
		{
			name:           "valid request method, valid content_type, but invalid url",
			reqMethod:      http.MethodPost,
			reqContentType: "text/plain",
			reqURL:         "some_url_without_protocol",
			respCode:       http.StatusBadRequest,
			isPositive:     false,
		},
		{
			name:           "valid request method, valid content_type, valid url -> 201 created",
			reqMethod:      http.MethodPost,
			reqContentType: "text/plain",
			reqURL:         "https://some_url_with_protocol",
			respCode:       http.StatusCreated,
			isPositive:     true,
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

			request, _ := http.NewRequest(test.reqMethod, "/", strings.NewReader(test.reqURL))
			request.Header.Add("Content-Type", test.reqContentType)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			api.Router.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			require.Equal(t, test.respCode, res.StatusCode)

			// проверяем, что урл сохранился в "БД"
			if test.isPositive {
				// получаем short_url из тела ответа
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)

				resURL := string(resBody)
				require.True(t, strings.HasPrefix(resURL, *config.BaseURL))

				splitted := strings.Split(resURL, "/")
				shortURL := splitted[len(splitted)-1]

				require.True(t, st.IsExists(shortURL))
			}
		})
	}
}
