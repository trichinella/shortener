package handler

import (
	"context"
	"net/http"
	"shortener/internal/app/repo"
	"time"
)

// PingDataBase что-то вроде healthcheck для БД
func PingDataBase(db repo.Pingable) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		childCtx, cancel := context.WithTimeout(r.Context(), time.Second*3)
		defer cancel()

		err := db.PingContext(childCtx)
		if err != nil {
			BadRequest(err, http.StatusInternalServerError)(w, r)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
