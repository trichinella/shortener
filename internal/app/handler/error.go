package handler

import (
	"net/http"
)

func BadRequest(err error, statusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err.Error(), statusCode)
	}
}
