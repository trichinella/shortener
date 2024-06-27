package repo

import (
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
)

// Repository Репозиторий с данными
type Repository interface {
	GetContraction(shortURL string) (*entity.Contraction, error)
	CreateContraction(originalURL string) *entity.Contraction
	HasContraction(shortURL string) bool
}

func GetRepo(cfg *config.MainConfig) Repository {
	if cfg.FileStoragePath == "" {
		return CreateMemoryRepository(cfg)
	}

	return CreateFileRepository(cfg)
}
