package repo

import (
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
)

// Repository Репозиторий с данными
type Repository interface {
	GetShortcut(shortURL string) (*entity.Shortcut, error)
	CreateShortcut(originalURL string) (*entity.Shortcut, error)
	HasShortcut(shortURL string) bool
}

func GetRepo() (Repository, error) {
	if config.State().FileStoragePath == "" {
		return CreateMemoryRepository(), nil
	}

	return CreateFileRepository()
}
