package handler

import (
	"fmt"
	"github.com/mailru/easyjson"
	"net/http"
	"shortener/internal/app/human"
	"shortener/internal/app/repo"
)

// CreateLinkPage Страница создания ссылки
func CreateLinkPage(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, ok := handleCreateLinkBody(w, r)
		if !ok {
			return
		}

		shortcut, err := repository.CreateShortcut(string(body))
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

// CreateLinkPageJSON Страница создания ссылки в формате JSON
func CreateLinkPageJSON(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, ok := handleCreateLinkBody(w, r)
		if !ok {
			return
		}

		inputURL := &InputURL{}
		err := easyjson.Unmarshal(body, inputURL)
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		if len(inputURL.URL) == 0 {
			BadRequest(fmt.Errorf("URL is empty"), http.StatusBadRequest)(w, r)
			return
		}

		shortcut, err := repository.CreateShortcut(inputURL.URL)
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		outputURL := &OutputURL{Result: human.GetFullShortURL(shortcut)}
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

type InputURL struct {
	URL string `json:"url"`
}

type OutputURL struct {
	Result string `json:"result"`
}
