package repo

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"shortener/internal/app/handler/inout"
)

// Repository Репозиторий с данными
type Repository interface {
	GetShortcutByShortURL(ctx context.Context, shortURL string) (*entity.Shortcut, error)
	GetShortcutByOriginalURL(ctx context.Context, originalURL string) (*entity.Shortcut, error)
	CreateShortcut(ctx context.Context, originalURL string) (*entity.Shortcut, error)
	CreateBatch(ctx context.Context, batchInput inout.ExternalBatchInput) (inout.ExternalBatchOutput, error)
	GetShortcutsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Shortcut, error)
}

// GetRepo выбор репозитория
// в приоритете postgres - если есть в конфиге запись
// далее файловый репозиторий
// ну а потом уже in memory
func GetRepo(db *sql.DB) (Repository, error) {
	if config.State().DatabaseDSN != "" {
		return CreatePostgresRepository(db), nil
	}

	if config.State().FileStoragePath != "" {
		return CreateFileRepository()
	}

	return CreateMemoryRepository(), nil
}

func HasShortcut(ctx context.Context, r Repository, shortURL string) bool {
	_, err := r.GetShortcutByShortURL(ctx, shortURL)

	return err == nil
}
