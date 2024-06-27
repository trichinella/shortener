package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"shortener/internal/app/repo"
	"testing"
)

func TestGetLinkPage(t *testing.T) {
	cfg := config.NewConfig()
	s := repo.CreateMemoryRepository(cfg)
	router := chi.NewRouter()
	router.Get(`/{shortUrl}`, GetLinkPage(s, cfg))
	ts := httptest.NewServer(router)
	defer ts.Close()

	type want struct {
		code     int
		url      string
		response string
	}
	tests := []struct {
		name        string
		contraction *entity.Contraction
		hash        string
		want        want
	}{
		{
			name:        "Base",
			contraction: s.CreateContraction("http://ya.ru"),
			want: want{
				code: http.StatusTemporaryRedirect,
				url:  "http://ya.ru",
			},
		},
		{
			name:        "Error",
			hash:        "itsnothabr",
			contraction: s.CreateContraction("http://ya.ru"),
			want: want{
				code:     http.StatusBadRequest,
				url:      "http://habr.ru",
				response: fmt.Errorf("unknown short url\n").Error(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shortUrl := test.contraction.ShortUrl
			if test.hash != "" {
				shortUrl = test.hash
			}

			req, err := http.NewRequest(http.MethodGet, string(ts.URL)+"/"+shortUrl, nil)
			req.Header.Set("Content-Type", "text/plain")

			require.NoError(t, err)

			client := ts.Client()
			hasRedirect := false
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				hasRedirect = true
				assert.Equal(t, test.want.url, req.URL.String())
				assert.Equal(t, test.want.code, req.Response.StatusCode)
				return http.ErrUseLastResponse
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

			if resp.StatusCode > 399 {
				// получаем и проверяем тело запроса
				assert.Equal(t, test.want.code, resp.StatusCode)
				assert.Equal(t, test.want.response, respBodyString)
				return
			}

			if !hasRedirect {
				t.Fatalf("There is not redirect")
			}
		})
	}
}
