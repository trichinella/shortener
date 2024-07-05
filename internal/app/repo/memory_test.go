package repo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"shortener/internal/app/handler/inout"
	"shortener/internal/app/human"
	"strings"
	"testing"
)

func TestMemoryRepository_CreateShortcut(t *testing.T) {
	tests := []struct {
		name string
		link string
	}{
		{
			name: "Пример #1",
			link: "https://ya.ru",
		},
		{
			name: "Пример #2",
			link: "https://lib.ru",
		},
		{
			name: "Пример #3",
			link: "https://ya.ru",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := CreateMemoryRepository()

			testShortcut, err := r.CreateShortcut(context.Background(), tt.link)
			require.NoError(t, err)

			if !strings.HasPrefix(human.GetFullShortURL(testShortcut), config.State().DisplayLink) {
				t.Errorf("Shortcut has incorrect prefix, got: %v,  want %v", human.GetFullShortURL(testShortcut), config.State().DisplayLink)
			}
		})
	}
}

func TestMemoryRepository_GetShortcut(t *testing.T) {
	type fields struct {
		Shortcuts map[string]entity.Shortcut
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
				Shortcuts: map[string]entity.Shortcut{
					"yaru12345": {
						ShortURL:    "yaru12345",
						OriginalURL: "http://ya.ru",
					},
					"qwerty": {
						ShortURL:    "qwerty",
						OriginalURL: "http://qwerty.ru",
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
				Shortcuts: map[string]entity.Shortcut{
					"yaru12345": {
						ShortURL:    "yaru12345",
						OriginalURL: "http://ya.ru",
					},
					"qwerty": {
						ShortURL:    "qwerty",
						OriginalURL: "http://qwerty.ru",
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
				Shortcuts: map[string]entity.Shortcut{
					"yaru12345": {
						ShortURL:    "yaru12345",
						OriginalURL: "http://ya.ru",
					},
					"qwerty": {
						ShortURL:    "qwerty",
						OriginalURL: "http://qwerty.ru",
					},
				},
			},
			hash:    "unknown",
			wantErr: fmt.Errorf("unknown short url"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MemoryRepository{
				Shortcuts: tt.fields.Shortcuts,
			}
			testShortcut, err := s.GetShortcut(context.Background(), tt.hash)

			if tt.wantErr == nil {
				require.NoError(t, tt.wantErr, err)
				assert.Equal(t, tt.want, testShortcut.OriginalURL)
			} else {
				require.Error(t, tt.wantErr, err)
			}
		})
	}
}

func TestMemoryRepository_HasShortcut(t *testing.T) {
	type fields struct {
		Shortcuts map[string]entity.Shortcut
		Host      string
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
				Shortcuts: map[string]entity.Shortcut{
					"yaru12345": {
						ShortURL:    "yaru12345",
						OriginalURL: "http://ya.ru",
					},
					"qwerty": {
						ShortURL:    "qwerty",
						OriginalURL: "http://qwerty.ru",
					},
				},
			},
			hash: "yaru12345",
			want: true,
		},
		{
			name: "There is not hash",
			fields: fields{
				Shortcuts: map[string]entity.Shortcut{
					"yaru12345": {
						ShortURL:    "yaru12345",
						OriginalURL: "http://ya.ru",
					},
					"qwerty": {
						ShortURL:    "qwerty",
						OriginalURL: "http://qwerty.ru",
					},
				},
			},
			hash: "unknown",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := MemoryRepository{
				Shortcuts: tt.fields.Shortcuts,
			}

			got := HasShortcut(context.Background(), &r, tt.hash)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMemoryRepository_CreateBatch(t *testing.T) {
	type fields struct {
		Shortcuts map[string]entity.Shortcut
	}
	type args struct {
		batchInput inout.ExternalBatchInput
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCount int
	}{
		{
			name: "1 в батче",
			fields: fields{
				map[string]entity.Shortcut{},
			},
			args: args{
				batchInput: []inout.ExternalInput{
					{ExternalID: "test", OriginalURL: "http://ya.ru"},
				},
			},
			wantCount: 1,
		},
		{
			name: "3 в батче",
			fields: fields{
				map[string]entity.Shortcut{},
			},
			args: args{
				batchInput: []inout.ExternalInput{
					{ExternalID: "test", OriginalURL: "http://ya.ru"},
					{ExternalID: "test2", OriginalURL: "http://ya1.ru"},
					{ExternalID: "test3", OriginalURL: "http://ya2.ru"},
				},
			},
			wantCount: 3,
		},
		{
			name: "0 в батче",
			fields: fields{
				map[string]entity.Shortcut{},
			},
			args: args{
				batchInput: []inout.ExternalInput{},
			},
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MemoryRepository{
				Shortcuts: tt.fields.Shortcuts,
			}
			gotResult, err := r.CreateBatch(context.Background(), tt.args.batchInput)

			require.Len(t, r.Shortcuts, tt.wantCount)
			require.Len(t, gotResult, tt.wantCount)
			require.NoError(t, err)

			for _, itemResult := range gotResult {
				has := false
				for _, itemInput := range tt.args.batchInput {
					if itemInput.ExternalID == itemResult.ExternalID {
						has = true
					}
				}
				require.True(t, has)
			}
		})
	}
}
