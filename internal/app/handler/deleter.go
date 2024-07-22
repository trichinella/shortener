package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"net/http"
	"shortener/internal/app/handler/inout"
	"shortener/internal/app/repo"
	"shortener/internal/app/service/authentification"
)

func DeleteUserURL(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(authentification.ContextUserID).(uuid.UUID)

		if !ok {
			BadRequest(fmt.Errorf("there is no user ID"), http.StatusUnauthorized)(w, r)
			return
		}

		body, ok := handleCreateLinkBody(w, r)
		if !ok {
			return
		}

		deletingList := inout.ShortURLList{}
		err := easyjson.Unmarshal(body, &deletingList)
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		err = repository.DeleteList(r.Context(), userID, deletingList)
		if err != nil {
			BadRequest(err, http.StatusInternalServerError)(w, r)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
