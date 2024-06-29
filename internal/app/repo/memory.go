package repo

import (
	"fmt"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"shortener/internal/app/random"
)

func CreateMemoryRepository(config *config.MainConfig) *MemoryRepository {
	return &MemoryRepository{
		Shortcuts: map[string]entity.Shortcut{},
		Config:    config,
	}
}

// MemoryRepository Основная структура
type MemoryRepository struct {
	Shortcuts map[string]entity.Shortcut
	Config    *config.MainConfig
}

// GetShortcut Получить сокращение на основе краткого URL
func (r *MemoryRepository) GetShortcut(shortURL string) (*entity.Shortcut, error) {
	shortcut, ok := r.Shortcuts[shortURL]

	if ok {
		return &shortcut, nil
	}

	return nil, fmt.Errorf("unknown short url")
}

func (r *MemoryRepository) HasShortcut(shortURL string) bool {
	_, err := r.GetShortcut(shortURL)

	return err == nil
}

// CreateShortcut Создать сокращение
func (r *MemoryRepository) CreateShortcut(originalURL string) (*entity.Shortcut, error) {
	var hash string
	for {
		hash = random.GenerateRandomString(7)
		if !r.HasShortcut(hash) {
			break
		}
	}

	shortcut := entity.Shortcut{
		OriginalURL: originalURL,
		ShortURL:    hash,
	}
	r.Shortcuts[shortcut.ShortURL] = shortcut

	return &shortcut, nil
}
