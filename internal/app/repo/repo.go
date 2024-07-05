package repo

import (
	"context"
	"database/sql"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"shortener/internal/app/handler/inout"
)

// Repository Репозиторий с данными
type Repository interface {
	GetShortcut(ctx context.Context, shortURL string) (*entity.Shortcut, error)
	CreateShortcut(ctx context.Context, originalURL string) (*entity.Shortcut, error)
	CreateBatch(ctx context.Context, batchInput inout.ExternalBatchInput) (inout.ExternalBatchOutput, error)
}

// GetRepo выбор репозитория
// в приоритете postgres - если есть в конфиге запись
// далее файловый репозиторий
// ну а потом уже in memory
func GetRepo(db *sql.DB) (Repository, error) {
	if config.State().DatabaseDSN != "" {
		return CreatePostgresRepository(db), nil
	}

	if config.State().FileStoragePath == "" {
		return CreateMemoryRepository(), nil
	}

	return CreateFileRepository()
}

func HasShortcut(ctx context.Context, r Repository, shortURL string) bool {
	_, err := r.GetShortcut(ctx, shortURL)

	return err == nil
}
