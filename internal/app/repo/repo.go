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

func GetRepo(cfg *config.MainConfig) (Repository, error) {
	if cfg.FileStoragePath == "" {
		return CreateMemoryRepository(cfg), nil
	}

	return CreateFileRepository(cfg)
}
