package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"shortener/internal/app/repo"
)

// GetLinkPage Страница получения ссылки
func GetLinkPage(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		val, err := repository.GetUserLink(hash)

		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		http.Redirect(w, r, val, http.StatusTemporaryRedirect)
	}
}
