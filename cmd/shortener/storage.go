package main

import (
	"fmt"
	"shortener/internal/app/random"
)

func CreateStore(baseLink string, port string) LocalRepository {
	return LocalRepository{
		BaseLink: baseLink + port + "/",
		Links:    map[string]string{},
	}
}

// Repository Задел на будущее (моки)
type Repository interface {
	GetLink(urlPath string) (string, error)
	CreateLink(link string) string
}

// Store Основная структура
type LocalRepository struct {
	Links    map[string]string
	BaseLink string
}

// Получить ссылку на основе URL
func (s LocalRepository) GetLink(urlPath string) (string, error) {
	val, ok := s.Links[urlPath[1:]]
	if !ok {
		return val, fmt.Errorf("unknown key")
	}

	return val, nil
}

// Создать ссылку - пока будем хранить в мапе
func (s LocalRepository) CreateLink(link string) string {
	var hash string
	for {
		hash = random.GenerateRandomString(7)
		if _, ok := s.Links[hash]; !ok {
			break
		}
	}

	s.Links[hash] = link

	return s.BaseLink + hash
}
