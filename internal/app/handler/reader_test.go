package handler

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"shortener/internal/app/entity"
	"shortener/internal/app/repo"
	"testing"
)

func TestGetShortcutPage(t *testing.T) {
	s := repo.CreateMemoryRepository()
	router := chi.NewRouter()
	router.Get(`/{shortURL}`, GetShortcutPage(s))
	ts := httptest.NewServer(router)
	defer ts.Close()

	type want struct {
		code     int
		url      string
		response string
	}
	tests := []struct {
		name        string
		shortcut    *entity.Shortcut
		shortURL    string
		originalURL string
		want        want
	}{
		{
			name:        "Base",
			originalURL: "http://ya.ru",
			want: want{
				code: http.StatusTemporaryRedirect,
				url:  "http://ya.ru",
			},
		},
		{
			name:        "Error",
			shortURL:    "itsnothabr",
			originalURL: "http://ya.ru",
			want: want{
				code:     http.StatusNotFound,
				url:      "http://habr.ru",
				response: fmt.Errorf("unknown short url\n").Error(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shortcut, err := s.CreateShortcut(context.Background(), test.originalURL)
			require.NoError(t, err)

			hash := shortcut.ShortURL
			if test.shortURL != "" {
				hash = test.shortURL
			}

			req, err := http.NewRequest(http.MethodGet, string(ts.URL)+"/"+hash, nil)
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
			require.NoError(t, err)

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

			require.True(t, hasRedirect, "There is not redirect")
		})
	}
}
