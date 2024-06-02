package main

import (
	"fmt"
	"shortener/internal/app/random"
)

func CreateLocalRepository(config *MainConfig) LocalRepository {
	return LocalRepository{
		Config:    config,
		UserLinks: map[string]string{},
	}
}

// Repository Задел на будущее (моки)
type Repository interface {
	GetUserLink(hash string) (string, error)
	CreateShortLink(userLink string) string
}

// LocalRepository Основная структура
type LocalRepository struct {
	UserLinks map[string]string
	Config    *MainConfig
}

// GetUserLink Получить ссылку на основе URL
func (s LocalRepository) GetUserLink(hash string) (string, error) {
	val, ok := s.UserLinks[hash]
	if !ok {
		return val, fmt.Errorf("unknown key")
	}

	return val, nil
}

// CreateShortLink Создать ссылку - пока будем хранить в мапе
func (s LocalRepository) CreateShortLink(userLink string) string {
	var hash string
	for {
		hash = random.GenerateRandomString(7)
		if _, ok := s.UserLinks[hash]; !ok {
			break
		}
	}

	s.UserLinks[hash] = userLink

	return s.Config.Protocol + "://" + s.Config.ShortLinkHost + "/" + hash
}
