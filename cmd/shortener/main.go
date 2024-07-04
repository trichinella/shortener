package main

import (
	"shortener/internal/app/repo"
	"shortener/internal/app/server"
)

func main() {
	logger := NewConsoleLogger()
	defer func() {
		_ = logger.Sync()
	}()

	db := repo.GetDB(logger)
	defer func() {
		err := db.Close()
		if err != nil {
			logger.Fatal(err.Error())
		}
	}()

	srv := server.CreateServer(logger, db)
	srv.Run()
}
