package handler

import (
	"database/sql"
	"net/http"
)

// PingDataBase что-то вроде healthcheck для БД
func PingDataBase(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			BadRequest(err, http.StatusInternalServerError)(w, r)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}
