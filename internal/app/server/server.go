package server

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"shortener/internal/app/config"
	"shortener/internal/app/handler"
	"shortener/internal/app/repo"
	"shortener/internal/app/server/middleware"
)

type CustomServer struct {
	Logger *zap.Logger
	Router *chi.Mux
	DB     *sql.DB
}

func (s *CustomServer) Run() {
	mainRepo, err := repo.GetRepo()
	if err != nil {
		panic(err)
	}

	s.Router = chi.NewRouter()
	s.Router.Use(middleware.Compress(s.Logger.Sugar()))
	s.Router.Use(middleware.LogMiddleware(s.Logger.Sugar()))
	fillHandler(s.Router, mainRepo, s.DB)

	s.Logger.Sugar().Infow("Listen and serve", "Host", config.State().ServerHost)
	err = http.ListenAndServe(config.State().ServerHost, s.Router)
	if err != nil {
		panic(err)
	}
}

func fillHandler(router chi.Router, repo repo.Repository, db *sql.DB) {
	router.Get(`/{shortURL}`, handler.GetLinkPage(repo))
	router.Post(`/api/shorten`, handler.CreateLinkPageJSON(repo))
	router.Post(`/`, handler.CreateLinkPage(repo))
	router.Get(`/ping`, handler.PingDataBase(db))
}

func CreateServer(logger *zap.Logger, db *sql.DB) CustomServer {
	return CustomServer{
		DB:     db,
		Logger: logger,
	}
}
