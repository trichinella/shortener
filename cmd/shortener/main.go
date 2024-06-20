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

	consoleLogger.Info("Create Server")
	srv := server.CreateServer(cfg, consoleLogger)
	consoleLogger.Info("Prepare Server")
	srv.Prepare()
	consoleLogger.Info("Start Server")
	srv.Start()
}
