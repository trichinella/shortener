package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestStore_CreateLinkPage(t *testing.T) {
	s := CreateLocalRepository(NewConfig())
	ts := httptest.NewServer(GetRouter(s))
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Base",
			url:  "http://ya.ru",
			want: want{
				code:        201,
				response:    s.Config.DisplayLink,
				contentType: "text/plain",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, get, _ := testRequest(t, ts, http.MethodPost, "/", strings.NewReader(test.url))
			err := res.Body.Close()
			if err != nil {
				require.NoError(t, err)
			}

			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			assert.Equal(t, true, strings.HasPrefix(get, test.want.response))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestStore_GetLinkPage(t *testing.T) {
	s := CreateLocalRepository(NewConfig())
	ts := httptest.NewServer(GetRouter(s))
	defer ts.Close()

	type want struct {
		code int
		url  string
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _, redirects := testRequest(t, ts, http.MethodGet, "/"+test.hash, nil)
			err := res.Body.Close()
			if err != nil {
				require.NoError(t, err)
			}

			if len(redirects) == 0 {
				t.Fatalf("There is not redirect")
			}

			assert.Equal(t, test.want.url, redirects[0].URL.String())
			assert.Equal(t, test.want.code, redirects[0].Code)
		})
	}
}

func TestBadRequest(t *testing.T) {
	type Wanted struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want Wanted
	}{
		{
			name: "Bad request",
			want: Wanted{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			BadRequest()(w, request)

			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer func() {
				err := res.Body.Close()
				if err != nil {
					require.NoError(t, err)
				}
			}()

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

type Redirect struct {
	URL  *url.URL
	Code int
}

func testRequest(t *testing.T, ts *httptest.Server, method string, path string, body io.Reader) (*http.Response, string, []Redirect) {
	req, err := http.NewRequest(method, ts.URL+path, body)
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

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			require.NoError(t, err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody), redirects
}
