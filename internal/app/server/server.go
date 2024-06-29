package server

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"shortener/internal/app/config"
	"shortener/internal/app/handler"
	"shortener/internal/app/repo"
	"shortener/internal/app/server/middleware"
)

type CustomServer struct {
	Config *config.MainConfig
	Logger *zap.Logger
	Router *chi.Mux
}

func (s *CustomServer) Run() {
	mainRepo, err := repo.GetRepo(s.Config)
	if err != nil {
		panic(err)
	}

	s.Router = chi.NewRouter()
	s.Router.Use(middleware.Compress(s.Logger.Sugar()))
	s.Router.Use(middleware.LogMiddleware(s.Logger.Sugar()))
	fillHandler(s.Router, mainRepo, s.Config)

	s.Logger.Sugar().Infow("Listen and serve", "Host", s.Config.ServerHost)
	err = http.ListenAndServe(s.Config.ServerHost, s.Router)
	if err != nil {
		panic(err)
	}
}

func fillHandler(router chi.Router, repo repo.Repository, cfg *config.MainConfig) {
	router.Get(`/{shortURL}`, handler.GetLinkPage(repo))
	router.Post(`/api/shorten`, handler.CreateLinkPageJSON(repo, cfg))
	router.Post(`/`, handler.CreateLinkPage(repo, cfg))
}

func CreateServer(config *config.MainConfig, logger *zap.Logger) CustomServer {
	return CustomServer{
		Config: config,
		Logger: logger,
	}
}
