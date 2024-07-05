package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"shortener/internal/app/handler"
	"shortener/internal/app/repo"
	"testing"
)

func TestCompress(t *testing.T) {
	Router := chi.NewRouter()
	Router.Use(Compress())
	repository := repo.CreateMemoryRepository()
	Router.Post(`/api/shorten`, handler.CreateShortcutJSON(repository))
	Router.Post(`/`, handler.CreateShortcutPlain(repository))
	srv := httptest.NewServer(Router)
	defer srv.Close()

	successBody := "body is empty\n"

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte{})
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", srv.URL+"/", nil)
		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Accept-Encoding", "")

		r.Header.Set("Content-Type", "text/html")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		defer func() {
			err := resp.Body.Close()
			require.NoError(t, err)
		}()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, successBody, string(b))
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString("")
		r := httptest.NewRequest("POST", srv.URL+"/api/shorten", buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")
		r.Header.Set("Content-Type", "text/html")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		defer func() {
			err := resp.Body.Close()
			require.NoError(t, err)
		}()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		b, err := io.ReadAll(zr)
		require.NoError(t, err)

		require.Equal(t, successBody, string(b))
	})
}
