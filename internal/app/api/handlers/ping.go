package handlers

import (
	"net/http"

	"github.com/Ppasha9/ya-shortener/internal/app/config"
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
)

func (h *handlers) PingHandler(w http.ResponseWriter, r *http.Request) {
	h.api.Logger.Info("Incoming GET pint request")

	if r.Method != http.MethodGet {
		// Принимаем только POST запросы
		h.api.Logger.Error("Invalid method", "method", r.Method)
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Пытаемся создать соединение с БД
	db, err := storage.NewDatabase(*config.DatabaseDSN)
	if err != nil {
		h.api.Logger.Error("Cannot create DB connection", "err", err.Error())
		http.Error(w, "Cannot create DB connection", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	h.api.Logger.Info("DB connection successfully created")
	w.WriteHeader(http.StatusOK)
}
