package repo

import (
	"fmt"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"shortener/internal/app/random"
)

func CreateMemoryRepository(config *config.MainConfig) *MemoryRepository {
	return &MemoryRepository{
		Contractions: []*entity.Contraction{},
		Config:       config,
	}
}

// MemoryRepository Основная структура
type MemoryRepository struct {
	Contractions []*entity.Contraction
	Config       *config.MainConfig
}

// GetContraction Получить сокращение на основе краткого URL
func (r *MemoryRepository) GetContraction(shortURL string) (*entity.Contraction, error) {
	for _, contraction := range r.Contractions {
		if contraction.ShortURL == shortURL {
			return contraction, nil
		}
	}

	return nil, fmt.Errorf("unknown short url")
}

func (r *MemoryRepository) HasContraction(shortURL string) bool {
	_, err := r.GetContraction(shortURL)

	return err == nil
}

// CreateContraction Создать сокращение
func (r *MemoryRepository) CreateContraction(originalURL string) *entity.Contraction {
	var hash string
	for {
		hash = random.GenerateRandomString(7)
		if !r.HasContraction(hash) {
			break
		}
	}

	contraction := &entity.Contraction{
		OriginalURL: originalURL,
		ShortURL:    hash,
	}
	r.Contractions = append(r.Contractions, contraction)

	return contraction
}