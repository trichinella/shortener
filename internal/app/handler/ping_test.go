package handler

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"shortener/internal/app/config"
	"shortener/internal/app/repo"
	"testing"
)

func TestPingDataBaseFailed(t *testing.T) {
	wantCode := http.StatusInternalServerError

	router := chi.NewRouter()
	logger, _ := zap.NewDevelopment()
	router.Get(`/ping`, PingDataBase(GetFakeDB(logger)))
	ts := httptest.NewServer(router)

	t.Run("Success ping", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, string(ts.URL)+"/ping", nil)
		req.Header.Set("Content-Type", "text/plain")

		require.NoError(t, err)

		client := ts.Client()
		resp, err := client.Do(req)
		require.NoError(t, err)

		defer func() {
			err := resp.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, wantCode, resp.StatusCode)
	})
}

func TestPingDataBase(t *testing.T) {
	wantCode := http.StatusOK

	router := chi.NewRouter()
	logger, _ := zap.NewDevelopment()
	router.Get(`/ping`, PingDataBase(repo.GetDB(logger)))
	ts := httptest.NewServer(router)

	t.Run("Success ping", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, string(ts.URL)+"/ping", nil)
		req.Header.Set("Content-Type", "text/plain")

		require.NoError(t, err)

		client := ts.Client()
		resp, err := client.Do(req)
		require.NoError(t, err)

		defer func() {
			err := resp.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, wantCode, resp.StatusCode)
	})
}

func GetFakeDB(logger *zap.Logger) *sql.DB {
	db, err := sql.Open("pgx", config.State().DatabaseDSN+"test")
	if err != nil {
		logger.Fatal(err.Error())
	}

	return db
}
