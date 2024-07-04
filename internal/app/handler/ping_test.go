package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	mocks "shortener/internal/mocks/repo"
	"testing"
)

func TestPingDataBase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Ping Pong (Success)",
			args: args{
				err: nil,
			},
			want: http.StatusOK,
		},
		{
			name: "Ping... (Error)",
			args: args{
				err: fmt.Errorf(""),
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		m := mocks.NewMockPingable(ctrl)
		m.EXPECT().Ping().Return(tt.args.err)

		router := chi.NewRouter()
		router.Get(`/ping`, PingDataBase(m))
		ts := httptest.NewServer(router)

		t.Run(tt.name, func(t *testing.T) {
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

			assert.Equal(t, tt.want, resp.StatusCode)
		})
	}
}
