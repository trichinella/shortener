package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"shortener/internal/app/config"
	"shortener/internal/app/repo"
)

// GetLinkPage Страница получения ссылки
func GetLinkPage(repository repo.Repository, cfg *config.MainConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortUrl := chi.URLParam(r, "shortUrl")
		contraction, err := repository.GetContraction(shortUrl)

		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		http.Redirect(w, r, contraction.OriginalUrl, http.StatusTemporaryRedirect)
	}
}
