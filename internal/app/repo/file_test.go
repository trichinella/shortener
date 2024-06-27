package repo

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"os/user"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"testing"
)

func TestFileRepository_init(t *testing.T) {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	type fields struct {
		Contractions []*entity.Contraction
		Config       *config.MainConfig
	}
	tests := []struct {
		name        string
		fields      fields
		fileContent string
		wantErr     bool
		want        int
	}{
		{
			name: "Пример из 5 сокращений",
			fields: fields{
				Contractions: []*entity.Contraction{},
				Config: &config.MainConfig{
					FileStoragePath: homeDir + "/f.json",
				},
			},
			fileContent: `{"uuid":0,"short_url":"lvs3iWf","original_url":"https://practicum.yandex.ru/"}
{"uuid":0,"short_url":"c7n4INA","original_url":"https://practicum.yandex.ru/"}
{"uuid":0,"short_url":"Fxuyi8z","original_url":"https://practicum.yandex.ru/"}
{"uuid":0,"short_url":"RuPCOGq","original_url":"https://practicum.yandex.ru/"}
{"uuid":0,"short_url":"1GgQeMp","original_url":"http://ya.ru"}`,
			wantErr: false,
			want:    5,
		},
		{
			name: "Пример из 0 сокращений",
			fields: fields{
				Contractions: []*entity.Contraction{},
				Config: &config.MainConfig{
					FileStoragePath: homeDir + "/f1.json",
				},
			},
			fileContent: ``,
			wantErr:     false,
			want:        0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.WriteFile(tt.fields.Config.FileStoragePath, []byte(tt.fileContent), 0666)
			if err != nil {
				t.Fatalf("Incorrect test: %v", err)
			}
			defer func(name string) {
				err := os.Remove(name)
				if err != nil {
					t.Fatalf("Incorrect test: %v", err)
				}
			}(tt.fields.Config.FileStoragePath)

			s := &FileRepository{
				Contractions: tt.fields.Contractions,
				Config:       tt.fields.Config,
			}

			if err := s.init(); (err != nil) != tt.wantErr {
				t.Errorf("init() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Len(t, s.Contractions, tt.want)
		})
	}
}

func TestFileRepository_CreateContraction(t *testing.T) {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	type fields struct {
		Contractions []*entity.Contraction
		Config       *config.MainConfig
	}
	type args struct {
		originalUrl string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *entity.Contraction
		wantCount int
	}{
		{
			name: "Пример из 1 сокращения",
			fields: fields{
				Contractions: []*entity.Contraction{},
				Config: &config.MainConfig{
					FileStoragePath: homeDir + "/q.json",
				},
			},
			args: struct {
				originalUrl string
			}{
				originalUrl: "http://ya.ru",
			},
			want: &entity.Contraction{
				OriginalUrl: "http://ya.ru",
			},
			wantCount: 1,
		},
		{
			name: "Пример из 3 сокращения",
			fields: fields{
				Contractions: []*entity.Contraction{},
				Config: &config.MainConfig{
					FileStoragePath: homeDir + "/q1.json",
				},
			},
			args: struct {
				originalUrl string
			}{
				originalUrl: "http://habr.ru",
			},
			want: &entity.Contraction{
				OriginalUrl: "http://habr.ru",
			},
			wantCount: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FileRepository{
				Contractions: tt.fields.Contractions,
				Config:       tt.fields.Config,
			}

			var gotContraction *entity.Contraction
			for i := 0; i < tt.wantCount; i++ {
				gotContraction = r.CreateContraction(tt.args.originalUrl)
			}

			defer func(name string) {
				err := os.Remove(name)
				if err != nil {
					t.Fatalf("Incorrect test: %v", err)
				}
			}(tt.fields.Config.FileStoragePath)
			file, err := os.OpenFile(tt.fields.Config.FileStoragePath, os.O_RDONLY, 0666)
			if err != nil {
				t.Fatalf("Incorrect test: %v", err)
			}
			cnt, _ := LineCounter(file)

			assert.Equal(t, tt.wantCount, cnt)
			assert.Equal(t, tt.want.OriginalUrl, gotContraction.OriginalUrl, "Original URL and got URL must be equal")
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
