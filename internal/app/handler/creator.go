package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"net/http"
	"shortener/internal/app/handler/inout"
	"shortener/internal/app/human"
	"shortener/internal/app/repo"
	"shortener/internal/app/service/authentification"
	"time"
)

// CreateShortcutPlain Страница создания ссылки
func CreateShortcutPlain(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkToken(r) {
			ctx, err := createToken(&w, r)
			if err != nil {
				BadRequest(err, http.StatusBadRequest)(w, r)
				return
			}
			r = r.WithContext(ctx)
		}

		body, ok := handleCreateLinkBody(w, r)
		if !ok {
			return
		}

		shortcut, err := repository.CreateShortcut(r.Context(), string(body))
		var duplicateShortcutErr *repo.DuplicateShortcutError
		if err != nil && !errors.As(err, &duplicateShortcutErr) {
			BadRequest(err, http.StatusInternalServerError)(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		if duplicateShortcutErr == nil {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}

		_, err = w.Write([]byte(human.GetFullShortURL(shortcut)))

		if err != nil {
			BadRequest(err, http.StatusInternalServerError)(w, r)
		}
	}
}

// CreateShortcutJSON Страница создания ссылки в формате JSON
func CreateShortcutJSON(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkToken(r) {
			ctx, err := createToken(&w, r)
			if err != nil {
				BadRequest(err, http.StatusBadRequest)(w, r)
				return
			}
			r = r.WithContext(ctx)
		}

		body, ok := handleCreateLinkBody(w, r)
		if !ok {
			return
		}

		inputURL := &inout.InputURL{}
		err := easyjson.Unmarshal(body, inputURL)
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		if len(inputURL.URL) == 0 {
			BadRequest(fmt.Errorf("URL is empty"), http.StatusBadRequest)(w, r)
			return
		}

		shortcut, err := repository.CreateShortcut(r.Context(), inputURL.URL)
		var duplicateShortcutErr *repo.DuplicateShortcutError
		if err != nil && !errors.As(err, &duplicateShortcutErr) {
			BadRequest(err, http.StatusInternalServerError)(w, r)
			return
		}

		outputURL := &inout.OutputURL{Result: human.GetFullShortURL(shortcut)}
		rawBytes, err := easyjson.Marshal(outputURL)
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if duplicateShortcutErr == nil {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
		_, err = w.Write(rawBytes)

		if err != nil {
			BadRequest(err, http.StatusInternalServerError)(w, r)
		}
	}
}

// CreateShortcutBatchJSON Страница создания ссылки батчем в формате JSON
func CreateShortcutBatchJSON(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkToken(r) {
			ctx, err := createToken(&w, r)
			if err != nil {
				BadRequest(err, http.StatusBadRequest)(w, r)
				return
			}
			r = r.WithContext(ctx)
		}

		body, ok := handleCreateLinkBody(w, r)
		if !ok {
			return
		}

		externalBatchInput := inout.ExternalBatchInput{}
		err := easyjson.Unmarshal(body, &externalBatchInput)
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		if len(externalBatchInput) == 0 {
			BadRequest(fmt.Errorf("batch is empty"), http.StatusBadRequest)(w, r)
			return
		}

		externalBatchOutput, err := repository.CreateBatch(r.Context(), externalBatchInput)
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		rawBytes, err := easyjson.Marshal(externalBatchOutput)
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(rawBytes)

		if err != nil {
			BadRequest(err, http.StatusInternalServerError)(w, r)
		}
	}
}

func handleCreateLinkBody(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	body, err := GetBody(r)
	if err != nil {
		BadRequest(err, http.StatusInternalServerError)(w, r)
		return nil, false
	}

	if len(body) == 0 {
		BadRequest(fmt.Errorf("body is empty"), http.StatusBadRequest)(w, r)
		return nil, false
	}

	return body, true
}

func checkToken(r *http.Request) bool {
	//есть ли кука?
	tokenCookie, err := r.Cookie("Authorization")

	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return false
	}

	if err != nil {
		return false
	}

	//если есть - получить claims из куки
	claims, err := authentification.GetClaims(tokenCookie.Value)
	if err != nil {
		return false
	}

	if claims.UserID == uuid.Nil {
		return false
	}

	return true
}

func createToken(w *http.ResponseWriter, r *http.Request) (context.Context, error) {
	signedToken, err := authentification.BuildJWTString()

	if err != nil {
		return nil, err
	}

	http.SetCookie(*w, &http.Cookie{
		Name:    "Authorization",
		Value:   signedToken,
		Expires: time.Now().Add(3 * time.Hour),
	})

	claims, err := authentification.GetClaims(signedToken)
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(r.Context(), authentification.ContextUserID, claims.UserID)

	return ctx, nil
}
