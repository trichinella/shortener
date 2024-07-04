package repo

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"shortener/internal/app/config"
)

func GetDB(logger *zap.Logger) *sql.DB {
	db, err := sql.Open("pgx", config.State().DatabaseDSN)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return db
}
