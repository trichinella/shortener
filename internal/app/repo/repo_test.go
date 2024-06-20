package repo

import (
	"fmt"
	"shortener/internal/app/config"
	"strings"
	"testing"
)

func TestStore_CreateShortLink(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name string
		args args
		link string
		want string
	}{
		{
			name: "Пример #1",
			args: args{
				host: "http://localhost:123",
			},
			link: "https://ya.ru",
			want: "http://localhost:123",
		},
		{
			name: "Пример #2",
			args: args{
				host: "http://example.site:443",
			},
			link: "https://lib.ru",
			want: "http://example.site:443",
		},
		{
			name: "Пример #3",
			args: args{
				host: "http://habr.ru:8080",
			},
			link: "https://ya.ru",
			want: "http://habr.ru:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig()
			cfg.DisplayLink = tt.args.host
			s := CreateLocalRepository(cfg)

			got := s.CreateShortLink(tt.link)

			if len(got) == 0 {
				t.Errorf("CreateShortLink() is empty, want > 0")
			}
			if !strings.HasPrefix(got, tt.want) {
				t.Errorf("CreateShortLink() has not prefix, got: %v,  want %v", got, tt.want)
			}
		})
	}
}

func TestStore_GetUserLink(t *testing.T) {
	type fields struct {
		UserLinks map[string]string
		Host      string
	}

	tests := []struct {
		name    string
		fields  fields
		hash    string
		want    string
		wantErr error
	}{
		{
			name: "Базовый функционал #1",
			fields: fields{
				UserLinks: map[string]string{
					"qwerty":    "http://qwerty.ru",
					"yaru12345": "http://ya.ru",
				},
			},
			hash:    "qwerty",
			want:    "http://qwerty.ru",
			wantErr: nil,
		},
		{
			name: "Базовый функционал #2",
			fields: fields{
				UserLinks: map[string]string{
					"qwerty":    "http://qwerty.ru",
					"yaru12345": "http://ya.ru",
				},
			},
			hash:    "yaru12345",
			want:    "http://ya.ru",
			wantErr: nil,
		},
		{
			name: "Не найдено #1",
			fields: fields{
				UserLinks: map[string]string{
					"qwerty":    "http://qwerty.ru",
					"yaru12345": "http://ya.ru",
				},
			},
			hash:    "unknown",
			wantErr: fmt.Errorf("unknown key"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig()
			s := LocalRepository{
				UserLinks: tt.fields.UserLinks,
				Config:    cfg,
			}
			got, err := s.GetUserLink(tt.hash)

			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetUserLink() got error = %v, want error = %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr == nil {
				t.Errorf("GetUserLink() got error = %v, want error = %v", err, tt.wantErr)
				return
			}

			if err == nil && tt.wantErr != nil {
				t.Errorf("GetUserLink() got error = %v, want error = %v", err, tt.wantErr)
				return
			}

			if err == nil && got != tt.want {
				t.Errorf("GetUserLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}
