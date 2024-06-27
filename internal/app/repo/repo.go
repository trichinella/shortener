package repo

import (
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
)

// Repository Репозиторий с данными
type Repository interface {
	GetContraction(shortUrl string) (*entity.Contraction, error)
	CreateContraction(originalUrl string) *entity.Contraction
	HasContraction(shortUrl string) bool
}

func GetRepo(cfg *config.MainConfig) Repository {
	//return CreateMemoryRepository(cfg)
	if "" == cfg.FileStoragePath {
		return CreateMemoryRepository(cfg)
	}
	return CreateFileRepository(cfg)
}
