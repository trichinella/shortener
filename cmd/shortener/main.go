package main

import (
	"shortener/internal/app/logging"
	"shortener/internal/app/repo"
	"shortener/internal/app/server"
)

func main() {
	defer func() {
		_ = logging.Logger.Sync()
	}()

	db := repo.GetDB()
	defer func() {
		err := db.Close()
		if err != nil {
			logging.Logger.Fatal(err.Error())
		}
	}()

	srv := server.CreateServer(db)
	srv.Run()
}
