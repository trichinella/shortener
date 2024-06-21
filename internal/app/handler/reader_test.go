package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"shortener/internal/app/config"
	"shortener/internal/app/repo"
	"strings"
	"testing"
)

func TestGetLinkPage(t *testing.T) {
	s := repo.CreateLocalRepository(config.NewConfig())
	router := chi.NewRouter()
	router.Get(`/{hash}`, GetLinkPage(s))
	ts := httptest.NewServer(router)
	defer ts.Close()

	type want struct {
		code     int
		url      string
		response string
	}
	tests := []struct {
		name string
		hash string
		want want
	}{
		{
			name: "Base",
			hash: strings.Split(s.CreateShortLink("http://ya.ru"), "/")[3],
			want: want{
				code: http.StatusTemporaryRedirect,
				url:  "http://ya.ru",
			},
		},
		{
			name: "Error",
			hash: "itsnothabr",
			want: want{
				code:     http.StatusBadRequest,
				url:      "http://habr.ru",
				response: fmt.Errorf("unknown key\n").Error(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, string(ts.URL)+"/"+test.hash, nil)
			req.Header.Set("Content-Type", "text/plain")

			require.NoError(t, err)

			client := ts.Client()
			var redirects []Redirect
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				redirects = append(redirects, Redirect{
					URL:  req.URL,
					Code: req.Response.StatusCode,
				})

				return nil
			}

			resp, err := client.Do(req)
			if err != nil {
				require.NoError(t, err)
			}

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			defer func() {
				err := resp.Body.Close()
				require.NoError(t, err)
			}()
			respBodyString := string(respBody)

			if resp.StatusCode > 299 {
				// получаем и проверяем тело запроса
				assert.Equal(t, test.want.code, resp.StatusCode)
				assert.Equal(t, test.want.response, respBodyString)
				return
			}

			if len(redirects) == 0 {
				t.Fatalf("There is not redirect")
			}

			assert.Equal(t, test.want.url, redirects[0].URL.String())
			assert.Equal(t, test.want.code, redirects[0].Code)
		})
	}
}

type Redirect struct {
	URL  *url.URL
	Code int
}
