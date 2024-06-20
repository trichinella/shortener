package main

import (
	"flag"
	"shortener/internal/app/config"
	"shortener/internal/app/server"
)

func main() {
	flag.Parse()

	//Заполнение конфига
	cfg := config.NewConfig()
	cfg.UpdateByOptions(config.BaseOptions)
	cfg.UpdateByEnv()

	consoleLogger := NewConsoleLogger()

	srv := server.CreateServer(cfg, consoleLogger)
	srv.Prepare()
	srv.Start()
}
