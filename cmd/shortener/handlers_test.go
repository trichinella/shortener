package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func CreateFakeRepo() FakeRepo {
	return FakeRepo{}
}

type FakeRepo struct {
}

// Создать ссылку - фейк
func (s FakeRepo) CreateLink(link string) string {
	return "faked_" + strings.Split(link, "//")[1] + "_hashed"
}

// Получить ссылку - фейк
func (s FakeRepo) GetLink(urlPath string) (string, error) {
	return "http://" + urlPath[1+len("faked_"):len(urlPath)-len("_hashed")], nil
}

func TestStore_CreateLinkPage(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		link string
		want want
	}{
		{
			name: "Base",
			link: "http://ya.ru",
			want: struct {
				code        int
				response    string
				contentType string
			}{code: 201, response: `faked_ya.ru_hashed`, contentType: "text/plain"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.link))
			// создаём новый Recorder
			w := httptest.NewRecorder()

			s := CreateFakeRepo()
			CreateLinkPage(s)(w, request)
			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestStore_GetLinkPage(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		target string
		want   want
	}{
		{
			name:   "Base",
			target: "/faked_ya.ru_hashed",
			want: struct {
				code        int
				response    string
				contentType string
			}{code: 307, response: "<a href=\"http://ya.ru\">Temporary Redirect</a>.\n\n", contentType: "text/html; charset=utf-8"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.target, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			s := CreateFakeRepo()
			GetLinkPage(s)(w, request)
			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_mainPage(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		target string
		method string
		want   want
	}{
		{
			name:   "Method PUT",
			target: "/",
			method: http.MethodPut,
			want: struct {
				code        int
				response    string
				contentType string
			}{code: 400, response: "", contentType: ""},
		},
		{
			name:   "Method DELETE",
			target: "/",
			method: http.MethodDelete,
			want: struct {
				code        int
				response    string
				contentType string
			}{code: 400, response: "", contentType: ""},
		},
		{
			name:   "Method HEAD",
			target: "/",
			method: http.MethodHead,
			want: struct {
				code        int
				response    string
				contentType string
			}{code: 400, response: "", contentType: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.target, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			s := CreateFakeRepo()
			mainPage(s)(w, request)
			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
