package repo

import (
	"context"
	"fmt"
	"shortener/internal/app/entity"
	"shortener/internal/app/handler/inout"
	"shortener/internal/app/human"
	"shortener/internal/app/random"
)

func CreateMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		Shortcuts: map[string]entity.Shortcut{},
	}
}

// MemoryRepository репозиторий на основе хранения в памяти
type MemoryRepository struct {
	Shortcuts map[string]entity.Shortcut
}

// GetShortcut Получить сокращение на основе краткого URL
func (r *MemoryRepository) GetShortcut(ctx context.Context, shortURL string) (*entity.Shortcut, error) {
	shortcut, ok := r.Shortcuts[shortURL]

	if ok {
		return &shortcut, nil
	}

	return nil, fmt.Errorf("unknown short url")
}

// CreateShortcut Создать сокращение
func (r *MemoryRepository) CreateShortcut(ctx context.Context, originalURL string) (*entity.Shortcut, error) {
	var hash string
	for {
		hash = random.GenerateRandomString(7)
		if !HasShortcut(ctx, r, hash) {
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

func (r *MemoryRepository) CreateBatch(ctx context.Context, batchInput inout.ExternalBatchInput) (result inout.ExternalBatchOutput, err error) {
	//нормальное поведение
	if len(batchInput) == 0 {
		return result, nil
	}

	for _, input := range batchInput {
		shortcut, err := r.CreateShortcut(ctx, input.OriginalURL)
		if err != nil {
			return nil, err
		}

		externalOutput := inout.ExternalOutput{}
		externalOutput.ExternalID = input.ExternalID
		externalOutput.ShortURL = human.GetFullShortURL(shortcut)

		result = append(result, externalOutput)
	}

	return result, nil
}
