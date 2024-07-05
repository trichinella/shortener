package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"shortener/internal/app/repo"
)

// GetShortcutPage Страница получения ссылки
func GetShortcutPage(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortcut, err := repository.GetShortcut(r.Context(), chi.URLParam(r, "shortURL"))

		if err != nil {
			BadRequest(err, http.StatusNotFound)(w, r)
			return
		}

		http.Redirect(w, r, shortcut.OriginalURL, http.StatusTemporaryRedirect)
	}
}
