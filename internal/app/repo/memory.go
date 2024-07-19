package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

// GetShortcutByShortURL Получить сокращение на основе краткого URL
func (r *MemoryRepository) GetShortcutByShortURL(ctx context.Context, shortURL string) (*entity.Shortcut, error) {
	shortcut, ok := r.Shortcuts[shortURL]

	if ok {
		return &shortcut, nil
	}

	return nil, fmt.Errorf("unknown short url")
}

// GetShortcutByOriginalURL Получить сокращение на основе краткого URL
func (r *MemoryRepository) GetShortcutByOriginalURL(ctx context.Context, originalURL string) (*entity.Shortcut, error) {
	for _, shortcut := range r.Shortcuts {
		if shortcut.OriginalURL == originalURL {
			return &shortcut, nil
		}
	}

	return nil, nil
}

// CreateShortcut Создать сокращение
func (r *MemoryRepository) CreateShortcut(ctx context.Context, originalURL string) (*entity.Shortcut, error) {
	shortcut, err := r.GetShortcutByOriginalURL(ctx, originalURL)

	if shortcut != nil {

		return shortcut, NewDuplicateShortcutError(err, shortcut)
	}

	var hash string
	for {
		hash = random.GenerateRandomString(7)
		if !HasShortcut(ctx, r, hash) {
			break
		}
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	userID, ok := ctx.Value("UserID").(uuid.UUID)
	if !ok {
		userID = uuid.Nil
	}

	shortcut = &entity.Shortcut{
		ID:          id,
		OriginalURL: originalURL,
		ShortURL:    hash,
		UserID:      userID,
	}
	r.Shortcuts[shortcut.ShortURL] = *shortcut

	return shortcut, nil
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

func (r *MemoryRepository) GetShortcutsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Shortcut, error) {
	var shortcuts []entity.Shortcut
	for _, shortcut := range r.Shortcuts {
		if shortcut.UserID == userID {
			shortcuts = append(shortcuts, shortcut)
		}
	}

	return shortcuts, nil
}
