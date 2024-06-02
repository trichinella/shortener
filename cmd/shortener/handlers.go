package main

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

// Руками делим запросы - либо получить ссылку, либо создать ее, иначе плохой запрос
func BadRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Страница создания ссылки
func CreateLinkPage(repository Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		if len(body) == 0 {
			BadRequest()(w, r)
			return
		}

		link := string(body)
		hashedLink := repository.CreateShortLink(link)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(hashedLink))

		if err != nil {
			panic(err)
		}
	}
}

// Страница получения ссылки
func GetLinkPage(repository Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		val, err := repository.GetUserLink(hash)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, val, http.StatusTemporaryRedirect)
	}
}
