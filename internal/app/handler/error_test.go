package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBadRequest(t *testing.T) {
	type Wanted struct {
		code int
	}
	tests := []struct {
		name string
		err  error
		code int
		want Wanted
	}{
		{
			name: "Bad request #1",
			err:  fmt.Errorf("ошибка"),
			code: http.StatusBadRequest,
			want: Wanted{
				code: 400,
			},
		},
		{
			name: "Bad request #2",
			err:  fmt.Errorf("нет доступа"),
			code: http.StatusForbidden,
			want: Wanted{
				code: 403,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			BadRequest(tt.err, tt.code)(w, request)

			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer func() {
				err := res.Body.Close()
				require.NoError(t, err)
			}()

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.err.Error()+"\n", string(resBody))
		})
	}
}
