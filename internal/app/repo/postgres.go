package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"shortener/internal/app/logging"
	"shortener/internal/app/random"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_postgresql.go -package=mocks shortener/internal/app/repo Pingable

func GetDB() *sql.DB {
	db, err := sql.Open("pgx", config.State().DatabaseDSN)
	if err != nil {
		logging.Logger.Fatal(err.Error())
	}

	return db
}

type Pingable interface {
	Ping() error
}

// PostgresRepository репозиторий на основе хранения в БД postgres
type PostgresRepository struct {
	DB *sql.DB
}

func CreatePostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		DB: db,
	}
}

// GetShortcut Получить сокращение на основе краткого URL
func (r *PostgresRepository) GetShortcut(ctx context.Context, shortURL string) (*entity.Shortcut, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	var shortcut entity.Shortcut
	row := r.DB.QueryRowContext(childCtx,
		"SELECT s.short_url, s.original_url, s.uuid FROM public.shortcuts s WHERE s.short_url = $1",
		shortURL)
	err := row.Scan(&shortcut.ShortURL, &shortcut.OriginalURL, &shortcut.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unknown short url")
		}

		return nil, err
	}

	return &shortcut, nil
}

// CreateShortcut Создать сокращение
// По-хорошему ее надо тестировать через тестовую базу. А как ее внедрить так, чтобы автотесты не упали - пока нет идей
func (r *PostgresRepository) CreateShortcut(ctx context.Context, originalURL string) (*entity.Shortcut, error) {
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

	shortcut := entity.Shortcut{
		ID:          id,
		OriginalURL: originalURL,
		ShortURL:    hash,
	}

	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	_, err = r.DB.ExecContext(childCtx,
		"INSERT INTO public.shortcuts (uuid, original_url, short_url)  VALUES ($1, $2, $3)",
		shortcut.ID, shortcut.OriginalURL, shortcut.ShortURL)

	if err != nil {
		return nil, err
	}

	return &shortcut, nil
}
