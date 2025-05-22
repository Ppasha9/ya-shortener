package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func WithLogging(h http.Handler, logger *slog.Logger) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Incomming request",
			"method", r.Method,
			"path", r.RequestURI,
		)

		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)
		duration := time.Since(start)

		logger.Info("Request response",
			"method", r.Method,
			"path", r.RequestURI,
			"duration", duration,
			"status", responseData.status,
			"size", responseData.size,
		)
	}

	return http.HandlerFunc(logFn)
}
