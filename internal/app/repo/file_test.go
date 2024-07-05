package repo

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"testing"
)

func TestFileRepository_init(t *testing.T) {
	fileStoragePath := config.State().FileStoragePath

	t.Setenv("FILE_STORAGE_PATH", os.TempDir()+"/f.json")
	t.Cleanup(func() {
		err := os.Setenv("FILE_STORAGE_PATH", fileStoragePath)
		if err != nil {
			require.NoError(t, err)
		}
	})

	type fields struct {
		Shortcuts map[string]entity.Shortcut
	}
	tests := []struct {
		name        string
		fields      fields
		fileContent string
		want        int
	}{
		{
			name: "Пример из 5 сокращений",
			fields: fields{
				Shortcuts: map[string]entity.Shortcut{},
			},
			fileContent: `{"uuid":"00000000-0000-0000-0000-000000000000","short_url":"lvs3iWf","original_url":"https://practicum.yandex.ru/"}
{"uuid":"00000000-0000-0000-0000-000000000000","short_url":"c7n4INA","original_url":"https://practicum.yandex.ru/"}
{"uuid":"00000000-0000-0000-0000-000000000000","short_url":"Fxuyi8z","original_url":"https://practicum.yandex.ru/"}
{"uuid":"00000000-0000-0000-0000-000000000000","short_url":"RuPCOGq","original_url":"https://practicum.yandex.ru/"}
{"uuid":"00000000-0000-0000-0000-000000000000","short_url":"1GgQeMp","original_url":"http://ya.ru"}`,
			want: 5,
		},
		{
			name: "Пример из 0 сокращений",
			fields: fields{
				Shortcuts: map[string]entity.Shortcut{},
			},
			fileContent: ``,
			want:        0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.WriteFile(config.State().FileStoragePath, []byte(tt.fileContent), 0666)
			require.NoError(t, err)

			defer func(name string) {
				err := os.Remove(name)
				require.NoError(t, err, "Incorrect test")
			}(config.State().FileStoragePath)

			s := &FileRepository{
				Shortcuts: tt.fields.Shortcuts,
			}

			err = s.init()
			require.NoError(t, err)

			assert.Len(t, s.Shortcuts, tt.want)
		})
	}
}

func TestFileRepository_CreateShortcut(t *testing.T) {
	fileStoragePath := config.State().FileStoragePath

	t.Setenv("FILE_STORAGE_PATH", os.TempDir()+"/q.json")
	t.Cleanup(func() {
		err := os.Setenv("FILE_STORAGE_PATH", fileStoragePath)
		if err != nil {
			require.NoError(t, err)
		}
	})

	var err error
	type fields struct {
		Shortcuts map[string]entity.Shortcut
	}
	type args struct {
		OriginalURL string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *entity.Shortcut
		wantCount int
	}{
		{
			name: "Пример из 1 сокращения",
			fields: fields{
				Shortcuts: map[string]entity.Shortcut{},
			},
			args: struct {
				OriginalURL string
			}{
				OriginalURL: "http://ya.ru",
			},
			want: &entity.Shortcut{
				OriginalURL: "http://ya.ru",
			},
			wantCount: 1,
		},
		{
			name: "Пример из 3 сокращения",
			fields: fields{
				Shortcuts: map[string]entity.Shortcut{},
			},
			args: struct {
				OriginalURL string
			}{
				OriginalURL: "http://habr.ru",
			},
			want: &entity.Shortcut{
				OriginalURL: "http://habr.ru",
			},
			wantCount: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FileRepository{
				Shortcuts: tt.fields.Shortcuts,
			}

			var gotShortcut *entity.Shortcut
			for i := 0; i < tt.wantCount; i++ {
				gotShortcut, err = r.CreateShortcut(context.Background(), tt.args.OriginalURL)
				require.NoError(t, err)
			}

			defer func(name string) {
				err := os.Remove(name)
				require.NoError(t, err)
			}(config.State().FileStoragePath)
			file, err := os.OpenFile(config.State().FileStoragePath, os.O_RDONLY, 0666)
			require.NoError(t, err)
			cnt, _ := LineCounter(file)

			assert.Equal(t, tt.wantCount, cnt)
			assert.Equal(t, tt.want.OriginalURL, gotShortcut.OriginalURL, "Original URL and got URL must be equal")
		})
	}
}

func LineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
