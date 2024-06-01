package main

import "testing"

func TestStore_createLink(t *testing.T) {
	type fields struct {
		Links    map[string]string
		BaseLink string
	}
	tests := []struct {
		name   string
		fields fields
		link   string
		want   string
	}{
		{
			name: "Пример при пустом хранилище",
			fields: fields{
				Links:    map[string]string{},
				BaseLink: "",
			},
			link: "https://ya.ru",
		},
		{
			name: "Пример при не пустом хранилище",
			fields: fields{
				Links:    map[string]string{},
				BaseLink: "",
			},
			link: "https://lib.ru",
		},
		{
			name: "Пример при повторной ссылке",
			fields: fields{
				Links:    map[string]string{},
				BaseLink: "",
			},
			link: "https://ya.ru",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := LocalRepository{
				Links:    tt.fields.Links,
				BaseLink: tt.fields.BaseLink,
			}
			if got := s.CreateLink(tt.link); len(got) == 0 {
				t.Errorf("createLink() is empty, want > 0")
			}
		})
	}
}

func TestStore_getLink(t *testing.T) {
	type fields struct {
		Links    map[string]string
		BaseLink string
	}

	tests := []struct {
		name    string
		fields  fields
		urlPath string
		want    string
		wantErr error
	}{
		{
			name: "Базовый функционал #1",
			fields: fields{
				Links: map[string]string{
					"qwerty":    "http://qwerty.ru",
					"yaru12345": "http://ya.ru",
				},
			},
			urlPath: "/qwerty",
			want:    "http://qwerty.ru",
			wantErr: nil,
		},
		{
			name: "Базовый функционал #2",
			fields: fields{
				Links: map[string]string{
					"qwerty":    "http://qwerty.ru",
					"yaru12345": "http://ya.ru",
				},
			},
			urlPath: "/yaru12345",
			want:    "http://ya.ru",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := LocalRepository{
				Links:    tt.fields.Links,
				BaseLink: tt.fields.BaseLink,
			}
			got, err := s.GetLink(tt.urlPath)
			if err != tt.wantErr {
				t.Errorf("getLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}
