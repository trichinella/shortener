package main

import (
	"io"
	"net/http"
)

// Руками делим запросы - либо получить ссылку, либо создать ее, иначе плохой запрос
func mainPage(repository Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetLinkPage(repository)(w, r)
			return
		}

		if r.Method == http.MethodPost {
			CreateLinkPage(repository)(w, r)
			return
		}

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

		link := string(body)
		hashedLink := repository.CreateLink(link)

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
		val, err := repository.GetLink(r.URL.Path)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, val, http.StatusTemporaryRedirect)
	}
}
