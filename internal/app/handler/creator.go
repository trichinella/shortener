package handler

import (
	"fmt"
	"github.com/mailru/easyjson"
	"net/http"
	"shortener/internal/app/handler/inout"
	"shortener/internal/app/human"
	"shortener/internal/app/repo"
)

// CreateShortcutPlain Страница создания ссылки
func CreateShortcutPlain(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, ok := handleCreateLinkBody(w, r)
		if !ok {
			return
		}

		shortcut, err := repository.CreateShortcut(r.Context(), string(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(human.GetFullShortURL(shortcut)))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// CreateShortcutJSON Страница создания ссылки в формате JSON
func CreateShortcutJSON(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		outputURL := &inout.OutputURL{Result: human.GetFullShortURL(shortcut)}
		rawBytes, err := easyjson.Marshal(outputURL)
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(rawBytes)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// CreateShortcutBatchJSON Страница создания ссылки батчем в формате JSON
func CreateShortcutBatchJSON(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleCreateLinkBody(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	body, err := GetBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, false
	}

	if len(body) == 0 {
		BadRequest(fmt.Errorf("body is empty"), http.StatusBadRequest)(w, r)
		return nil, false
	}

	return body, true
}
