package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type CustomServer struct {
	Config *MainConfig
}

func (s CustomServer) Start() {
	repo := CreateLocalRepository(s.Config)
	router := GetRouter(repo)

	err := http.ListenAndServe(s.Config.ServerHost, router)
	if err != nil {
		panic(err)
	}
}

func GetRouter(repo Repository) chi.Router {
	r := chi.NewRouter()

	r.Get(`/{hash}`, GetLinkPage(repo))
	r.Post(`/`, CreateLinkPage(repo))

	return r
}

func CreateServer(config *MainConfig) CustomServer {
	return CustomServer{
		Config: config,
	}
}
