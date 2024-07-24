package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"net/http"
	"shortener/internal/app/repo"
	"shortener/internal/app/service/authentification"
	"shortener/internal/app/service/mapper"
)

// GetShortcutPage Страница получения ссылки
func GetShortcutPage(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortcut, err := repository.GetShortcutByShortURL(r.Context(), chi.URLParam(r, "shortURL"))

		if err != nil {
			BadRequest(err, http.StatusNotFound)(w, r)
			return
		}

		if shortcut.DeletedDate != nil {
			BadRequest(fmt.Errorf("shortcut is unavailable"), http.StatusGone)(w, r)
			return
		}

		http.Redirect(w, r, shortcut.OriginalURL, http.StatusTemporaryRedirect)
	}
}

// GetShortcutsByUser Страница получения ссылок пользователя
func GetShortcutsByUser(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(authentification.ContextUserID).(uuid.UUID)

		if !ok {
			BadRequest(fmt.Errorf("there is no user ID"), http.StatusUnauthorized)(w, r)
			return
		}

		shortcuts, err := repository.GetShortcutsByUserID(r.Context(), userID)
		if err != nil {
			BadRequest(err, http.StatusNotFound)(w, r)
			return
		}

		if len(shortcuts) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		rawBytes, err := easyjson.Marshal(mapper.GetBaseShortcutListFromShortcuts(shortcuts))
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(rawBytes)

		if err != nil {
			BadRequest(err, http.StatusInternalServerError)(w, r)
		}

	}
}
