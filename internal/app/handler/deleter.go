package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"net/http"
	"shortener/internal/app/handler/inout"
	"shortener/internal/app/logging"
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

		go func() {
			err = repository.DeleteList(context.Background(), userID, deletingList)
			if err != nil {
				logging.Sugar.Error(fmt.Errorf("при удалении списка возникла ошибка %w", err))
			}
		}()

		w.WriteHeader(http.StatusAccepted)
	}
}
