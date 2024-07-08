package handler

import (
	"net/http"
	"shortener/internal/app/repo"
)

// PingDataBase что-то вроде healthcheck для БД
func PingDataBase(db repo.Pingable) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			BadRequest(err, http.StatusInternalServerError)(w, r)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
