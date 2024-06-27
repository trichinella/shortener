package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"shortener/internal/app/repo"
)

// GetLinkPage Страница получения ссылки
func GetLinkPage(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contraction, err := repository.GetContraction(chi.URLParam(r, "shortURL"))

		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		http.Redirect(w, r, contraction.OriginalURL, http.StatusTemporaryRedirect)
	}
}
