package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"shortener/internal/app/repo"
	"strings"
	"testing"
)

func TestCreateShortcutPlain(t *testing.T) {
	s := repo.CreateMemoryRepository()
	router := chi.NewRouter()
	router.Post(`/`, CreateShortcutPlain(s))
	ts := httptest.NewServer(router)
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name        string
		body        io.Reader
		contentType string
		want        want
	}{
		{
			name:        "Работоспособный вариант в пустое хранилище",
			body:        strings.NewReader("http://ya.ru"),
			contentType: "text/plain",
			want: want{
				code:        201,
				contentType: "text/plain",
			},
		},
		{
			name:        "Ту же самую ссылку повторно отправить",
			body:        strings.NewReader("http://ya.ru"),
			contentType: "text/plain",
			want: want{
				code:        409,
				contentType: "text/plain",
			},
		},
		{
			name:        "Пустое тело запроса проверить",
			body:        strings.NewReader(""),
			contentType: "text/plain",
			want: want{
				code:        400,
				response:    fmt.Errorf("body is empty\n").Error(),
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, ts.URL, test.body)
			require.NoError(t, err)
			req.Header.Set("Content-Type", test.contentType)

			client := ts.Client()
			resp, err := client.Do(req)
			require.NoError(t, err)

			defer func() {
				err := resp.Body.Close()
				require.NoError(t, err)
			}()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			respBodyString := string(respBody)

			// проверяем код ответа
			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))

			if (resp.StatusCode >= 200 && resp.StatusCode <= 299) || resp.StatusCode == 409 {
				// получаем и проверяем тело запроса
				assert.Equal(t, true, strings.HasPrefix(respBodyString, test.want.response))
			} else {
				assert.Equal(t, test.want.response, respBodyString)
			}
		})
	}
}

func TestCreateShortcutJSON(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		target string
		body   string
		want   want
	}{
		{
			name:   "Success",
			target: "/api/shorten",
			body:   "{\"url\":\"http://ya.ru\"}",
			want: want{
				code:        201,
				contentType: "application/json",
			},
		},
		{
			name:   "Same request",
			target: "/api/shorten",
			body:   "{\"url\":\"http://ya.ru\"}",
			want: want{
				code:        409,
				contentType: "application/json",
			},
		},
		{
			name:   "Empty Body",
			target: "/api/shorten",
			body:   "",
			want: want{
				code:     400,
				response: fmt.Errorf("body is empty\n").Error(),
			},
		},
		{
			name:   "Incorrect JSON",
			target: "/api/shorten",
			body:   "{}}",
			want: want{
				code: 400,
			},
		},
		{
			name:   "Empty url",
			target: "/api/shorten",
			body:   "{\"url\":\"\"}",
			want: want{
				code:     400,
				response: "URL is empty\n",
			},
		},
		{
			name:   "Incorrect param",
			target: "/api/shorten",
			body:   "{\"url2\":\"http://ya.ru\"}",
			want: want{
				code: 400,
			},
		},
		{
			name:   "Incorrect param register",
			target: "/api/shorten",
			body:   "{\"URL\":\"http://ya.ru\"}",
			want: want{
				code: 400,
			},
		},
	}
	s := repo.CreateMemoryRepository()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.target, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			CreateShortcutJSON(s)(w, req)
			res := w.Result()

			defer func() {
				err := res.Body.Close()
				require.NoError(t, err)
			}()

			data, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, res.StatusCode)

			if res.StatusCode >= 200 && res.StatusCode <= 299 {
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
				assert.Regexpf(t, regexp.MustCompile(`^{"result":"http://localhost:8080/\S+?"}$`), string(data), "Incorrect response")
			} else {
				if tt.want.response != "" {
					assert.Equal(t, tt.want.response, string(data))
				}
			}
		})
	}
}
