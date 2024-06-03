package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

// BadRequest Руками делим запросы - либо получить ссылку, либо создать ее, иначе плохой запрос
func BadRequest(err error, statusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err.Error(), statusCode)
	}
}

// CreateLinkPage Страница создания ссылки
func CreateLinkPage(repository Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		err = r.Body.Close()
		if err != nil {
			panic(err)
		}

		if len(body) == 0 {
			BadRequest(fmt.Errorf("body is empty"), http.StatusBadRequest)(w, r)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "text/plain" {
			BadRequest(fmt.Errorf("Content-Type must be \"text/plain\""), http.StatusBadRequest)(w, r)
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

// GetLinkPage Страница получения ссылки
func GetLinkPage(repository Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		val, err := repository.GetUserLink(hash)

		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		http.Redirect(w, r, val, http.StatusTemporaryRedirect)
	}
}
