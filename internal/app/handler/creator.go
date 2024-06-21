package handler

import (
	"fmt"
	"github.com/mailru/easyjson"
	"io"
	"net/http"
	"shortener/internal/app/repo"
)

// CreateLinkPage Страница создания ссылки
func CreateLinkPage(repository repo.Repository) http.HandlerFunc {
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

// CreateLinkPageJson Похожа на CreateLinkPage, но отличается тем, что есть JSON. Думал об объединении методов
// пришел к выводу, что не рентабельно
func CreateLinkPageJson(repository repo.Repository) http.HandlerFunc {
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

		inputUrl := &InputUrl{}
		err = easyjson.Unmarshal(body, inputUrl)
		if err != nil {
			BadRequest(err, http.StatusBadRequest)(w, r)
			return
		}

		if len(inputUrl.Url) == 0 {
			BadRequest(fmt.Errorf("URL is empty"), http.StatusBadRequest)(w, r)
			return
		}

		hashedLink := repository.CreateShortLink(inputUrl.Url)

		outputUrl := &OutputUrl{Result: hashedLink}
		rawBytes, err := easyjson.Marshal(outputUrl)
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

type InputUrl struct {
	Url string `json:"url"`
}

type OutputUrl struct {
	Result string `json:"result"`
}
