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

func (s *CustomServer) Prepare() {
	mainRepo := repo.GetRepo(s.Config)
	s.Router = chi.NewRouter()
	s.Router.Use(middleware.Compress(s.Logger.Sugar()))
	s.Router.Use(middleware.LogMiddleware(s.Logger.Sugar()))
	fillHandler(s.Router, mainRepo, s.Config)
}

func (s *CustomServer) Start() {
	s.Logger.Sugar().Infow("Listen and serve", "Host", s.Config.ServerHost)
	err := http.ListenAndServe(s.Config.ServerHost, s.Router)
	if err != nil {
		panic(err)
	}
}

func fillHandler(router chi.Router, repo repo.Repository, cfg *config.MainConfig) {
	router.Get(`/{shortUrl}`, handler.GetLinkPage(repo, cfg))
	router.Post(`/api/shorten`, handler.CreateLinkPageJSON(repo, cfg))
	router.Post(`/`, handler.CreateLinkPage(repo, cfg))
}

func CreateServer(config *config.MainConfig, logger *zap.Logger) CustomServer {
	return CustomServer{
		Config: config,
		Logger: logger,
	}
}
