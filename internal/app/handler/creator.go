package handler

import (
	"fmt"
	"github.com/mailru/easyjson"
	"net/http"
	"shortener/internal/app/repo"
)

// CreateLinkPage Страница создания ссылки
func CreateLinkPage(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := GetBody(r)
		if err != nil {
			panic(err)
		}

		err = r.Body.Close()
		if err != nil {
			panic(err)
		}

		if len(string(body)) == 0 {
			BadRequest(fmt.Errorf("body is empty"), http.StatusBadRequest)(w, r)
			return
		}

		hashedLink := repository.CreateShortLink(string(body))

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(hashedLink))

		if err != nil {
			panic(err)
		}
	}
}

// CreateLinkPageJSON Похожа на CreateLinkPage, но отличается тем, что есть JSON. Думал об объединении методов
// пришел к выводу, что не рентабельно
func CreateLinkPageJSON(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := GetBody(r)
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

		inputURL := &InputURL{}
		err = easyjson.Unmarshal(body, inputURL)
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		if len(inputURL.URL) == 0 {
			BadRequest(fmt.Errorf("URL is empty"), http.StatusBadRequest)(w, r)
			return
		}

		hashedLink := repository.CreateShortLink(inputURL.URL)

		outputURL := &OutputURL{Result: hashedLink}
		rawBytes, err := easyjson.Marshal(outputURL)
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(rawBytes)

		if err != nil {
			panic(err)
		}
	}
}

type InputURL struct {
	URL string `json:"url"`
}

type OutputURL struct {
	Result string `json:"result"`
}
