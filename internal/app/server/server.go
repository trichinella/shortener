package server

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"net/http"
	"shortener/internal/app/config"
	"shortener/internal/app/handler"
	"shortener/internal/app/logging"
	"shortener/internal/app/repo"
	"shortener/internal/app/server/middleware"
)

type CustomServer struct {
	Router *chi.Mux
	DB     *sql.DB
}

func (s *CustomServer) Run() {
	mainRepo, err := repo.GetRepo(s.DB)

	if err != nil {
		logging.Sugar.Fatal(err)
	}

	s.Router = chi.NewRouter()
	s.Router.Use(middleware.AuthMiddleware())
	s.Router.Use(middleware.Compress())
	s.Router.Use(middleware.LogMiddleware())
	fillHandler(s.Router, mainRepo, s.DB)

	logging.Sugar.Infow("Listen and serve", "Host", config.State().ServerHost)
	err = http.ListenAndServe(config.State().ServerHost, s.Router)
	if err != nil {
		logging.Sugar.Fatal(err)
	}
}

func fillHandler(router chi.Router, repo repo.Repository, db *sql.DB) {
	router.Get(`/{shortURL}`, handler.GetShortcutPage(repo))
	router.Post(`/api/shorten`, handler.CreateShortcutJSON(repo))
	router.Post(`/api/shorten/batch`, handler.CreateShortcutBatchJSON(repo))
	router.Get(`/api/user/urls`, handler.GetShortcutsByUser(repo))
	router.Delete(`/api/user/urls`, handler.DeleteUserURL(repo))
	router.Post(`/`, handler.CreateShortcutPlain(repo))
	router.Get(`/ping`, handler.PingDataBase(db))
}

func CreateServer(db *sql.DB) CustomServer {
	return CustomServer{
		DB: db,
	}
}
