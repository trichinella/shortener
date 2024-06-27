package repo

import (
	"fmt"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"shortener/internal/app/human"
	"strings"
	"testing"
)

func TestMemoryRepository_CreateContraction(t *testing.T) {
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
			r := CreateMemoryRepository(cfg)

			testContraction := r.CreateContraction(tt.link)

			if !strings.HasPrefix(human.GetFullShortUrl(cfg, testContraction), tt.want) {
				t.Errorf("Contraction has incorrect prefix, got: %v,  want %v", human.GetFullShortUrl(cfg, testContraction), tt.want)
			}
		})
	}
}

func TestMemoryRepository_GetContraction(t *testing.T) {
	type fields struct {
		Contractions []*entity.Contraction
		Host         string
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
				Contractions: []*entity.Contraction{
					{
						ShortUrl:    "yaru12345",
						OriginalUrl: "http://ya.ru",
					},
					{
						ShortUrl:    "qwerty",
						OriginalUrl: "http://qwerty.ru",
					},
				},
			},
			hash:    "qwerty",
			want:    "http://qwerty.ru",
			wantErr: nil,
		},
		{
			name: "Базовый функционал #2",
			fields: fields{
				Contractions: []*entity.Contraction{
					{
						ShortUrl:    "yaru12345",
						OriginalUrl: "http://ya.ru",
					},
					{
						ShortUrl:    "qwerty",
						OriginalUrl: "http://qwerty.ru",
					},
				},
			},
			hash:    "yaru12345",
			want:    "http://ya.ru",
			wantErr: nil,
		},
		{
			name: "Не найдено #1",
			fields: fields{
				Contractions: []*entity.Contraction{
					{
						ShortUrl:    "yaru12345",
						OriginalUrl: "http://ya.ru",
					},
					{
						ShortUrl:    "qwerty",
						OriginalUrl: "http://qwerty.ru",
					},
				},
			},
			hash:    "unknown",
			wantErr: fmt.Errorf("unknown short url"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig()
			s := MemoryRepository{
				Contractions: tt.fields.Contractions,
				Config:       cfg,
			}
			testContraction, err := s.GetContraction(tt.hash)

			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetContraction() got error = %v, want error = %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr == nil {
				t.Errorf("GetContraction() got error = %v, want error = %v", err, tt.wantErr)
				return
			}

			if err == nil && tt.wantErr != nil {
				t.Errorf("GetContraction() got error = %v, want error = %v", err, tt.wantErr)
				return
			}

			if err == nil && testContraction.OriginalUrl != tt.want {
				t.Errorf("GetContraction() got = %v, want %v", testContraction.OriginalUrl, tt.want)
			}
		})
	}
}

func TestMemoryRepository_HasContraction(t *testing.T) {
	type fields struct {
		Contractions []*entity.Contraction
		Host         string
	}

	tests := []struct {
		name   string
		fields fields
		hash   string
		want   bool
	}{
		{
			name: "There is hash",
			fields: fields{
				Contractions: []*entity.Contraction{
					{
						ShortUrl:    "yaru12345",
						OriginalUrl: "http://ya.ru",
					},
					{
						ShortUrl:    "qwerty",
						OriginalUrl: "http://qwerty.ru",
					},
				},
			},
			hash: "yaru12345",
			want: true,
		},
		{
			name: "There is not hash",
			fields: fields{
				Contractions: []*entity.Contraction{
					{
						ShortUrl:    "yaru12345",
						OriginalUrl: "http://ya.ru",
					},
					{
						ShortUrl:    "qwerty",
						OriginalUrl: "http://qwerty.ru",
					},
				},
			},
			hash: "unknown",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MemoryRepository{
				Contractions: tt.fields.Contractions,
			}
			if got := s.HasContraction(tt.hash); got != tt.want {
				t.Errorf("HasContraction() = %v, want %v", got, tt.want)
			}
		})
	}
}
