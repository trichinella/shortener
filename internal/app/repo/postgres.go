package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"shortener/internal/app/handler/inout"
	"shortener/internal/app/human"
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
	execMigrations()
	return &PostgresRepository{
		DB: db,
	}
}

func execMigrations() {
	m, err := migrate.New(
		"file://internal/migrations",
		config.State().DatabaseDSN)
	if err != nil {
		logging.Sugar.Fatal(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logging.Sugar.Fatal(err)
	}
}

// GetShortcutByShortURL Получить сокращение на основе краткого URL
func (r *PostgresRepository) GetShortcutByShortURL(ctx context.Context, shortURL string) (*entity.Shortcut, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	var shortcut entity.Shortcut
	row := r.DB.QueryRowContext(childCtx,
		"SELECT s.uuid, s.original_url, s.short_url, s.created_date FROM public.shortcuts s WHERE s.short_url = $1",
		shortURL)
	err := row.Scan(&shortcut.ID, &shortcut.OriginalURL, &shortcut.ShortURL, &shortcut.CreatedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unknown short url")
		}

		return nil, err
	}

	return &shortcut, nil
}

// GetShortcutByOriginalURL Получить сокращение на основе оригинального URL
func (r *PostgresRepository) GetShortcutByOriginalURL(ctx context.Context, originalURL string) (*entity.Shortcut, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	var shortcut entity.Shortcut
	row := r.DB.QueryRowContext(childCtx,
		"SELECT s.uuid, s.original_url, s.short_url, s.created_date FROM public.shortcuts s WHERE s.original_url = $1",
		originalURL)
	err := row.Scan(&shortcut.ID, &shortcut.OriginalURL, &shortcut.ShortURL, &shortcut.CreatedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			//это штатная ситуация
			return nil, nil
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

	row := r.DB.QueryRowContext(childCtx,
		`INSERT INTO public.shortcuts (uuid, original_url, short_url, user_id) VALUES ($1, $2, $3, $4)
	returning uuid,
	original_url,
	short_url,
	created_date
`,
		shortcut.ID, shortcut.OriginalURL, shortcut.ShortURL, ctx.Value("UserID"))

	err = row.Scan(&shortcut.ID, &shortcut.OriginalURL, &shortcut.ShortURL, &shortcut.CreatedDate)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			shortcut, errInternal := r.GetShortcutByOriginalURL(ctx, shortcut.OriginalURL)
			if errInternal != nil {
				return nil, errInternal
			}

			return shortcut, NewDuplicateShortcutError(err, shortcut)
		}
	}

	if err != nil {
		return nil, err
	}

	return &shortcut, nil
}

func (r *PostgresRepository) CreateBatch(ctx context.Context, batchInput inout.ExternalBatchInput) (result inout.ExternalBatchOutput, err error) {
	//нормальное поведение
	if len(batchInput) == 0 {
		return result, nil
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, externalInput := range batchInput {
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
			OriginalURL: externalInput.OriginalURL,
			ShortURL:    hash,
		}

		_, err = tx.ExecContext(ctx,
			"INSERT INTO public.shortcuts (uuid, original_url, short_url, user_id) VALUES ($1, $2, $3, $4)",
			shortcut.ID, shortcut.OriginalURL, shortcut.ShortURL, ctx.Value("UserID"))

		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return nil, errRollback
			}
			return nil, err
		}

		externalOutput := inout.ExternalOutput{}
		externalOutput.ExternalID = externalInput.ExternalID
		externalOutput.ShortURL = human.GetFullShortURL(&shortcut)

		result = append(result, externalOutput)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PostgresRepository) GetShortcutsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Shortcut, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	var shortcuts []entity.Shortcut
	rows, err := r.DB.QueryContext(childCtx,
		"SELECT s.uuid, s.original_url, s.short_url, s.created_date FROM public.shortcuts s WHERE s.user_id = $1",
		userID)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logging.Sugar.Fatalf("Can not close rows: %s", err.Error())
		}
	}(rows)

	// пробегаем по всем записям
	for rows.Next() {
		var shortcut entity.Shortcut
		err := rows.Scan(&shortcut.ID, &shortcut.OriginalURL, &shortcut.ShortURL, &shortcut.CreatedDate)
		if err != nil {
			return nil, err
		}

		shortcuts = append(shortcuts, shortcut)
	}

	return shortcuts, nil
}
