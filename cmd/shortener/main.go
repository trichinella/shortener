package main

import (
	"io"
	"net/http"
	"shortener/internal/app/random"
)

type Store struct {
	Links    map[string]string
	BaseLink string
}

func (s Store) createLink(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	link := string(body)

	var hash string
	for {
		hash = random.GenerateRandomString(7)
		if _, ok := s.Links[hash]; !ok {
			break
		}
	}

	s.Links[hash] = link
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(s.BaseLink + hash))

	if err != nil {
		panic(err)
	}
}

func (s Store) getLink(w http.ResponseWriter, r *http.Request) {
	val, ok := s.Links[r.URL.Path[1:]]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, val, http.StatusTemporaryRedirect)
}

func (s Store) mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.getLink(w, r)
		return
	}

	if r.Method == http.MethodPost {
		s.createLink(w, r)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func main() {
	mux := http.NewServeMux()
	port := "8080"

	s := &Store{
		BaseLink: "http://localhost:" + port + "/",
		Links:    map[string]string{},
	}
	mux.HandleFunc(`/`, s.mainPage)
	err := http.ListenAndServe(`:`+port, mux)
	if err != nil {
		panic(err)
	}
}
