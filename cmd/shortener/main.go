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

	logger := NewConsoleLogger()
	defer func() {
		_ = logger.Sync()
	}()

	srv := server.CreateServer(cfg, logger)
	srv.Run()
}
