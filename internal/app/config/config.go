package config

import (
	"os"
	"strings"
)

type MainConfig struct {
	DisplayLink     string
	ServerHost      string
	FileStoragePath string
	DatabaseDSN     string
}

func (config *MainConfig) updateByOptions(o options) {
	config.ServerHost = o.ServerHost
	config.DisplayLink = strings.Trim(o.DisplayLink, "/")
	config.FileStoragePath = o.FileStoragePath
	config.DatabaseDSN = o.DatabaseDSN
}

func (config *MainConfig) updateByEnv() {
	envDisplayLink := os.Getenv("BASE_URL")
	if envDisplayLink != "" {
		config.DisplayLink = envDisplayLink
	}

	envServerAddr := os.Getenv("SERVER_ADDRESS")
	if envServerAddr != "" {
		config.ServerHost = envServerAddr
	}

	envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if envFileStoragePath != "" {
		config.FileStoragePath = envFileStoragePath
	}

	envDatabaseDSN := os.Getenv("DATABASE_DSN")
	if envDatabaseDSN != "" {
		config.DatabaseDSN = envDatabaseDSN
	}
}

func newConfig() *MainConfig {
	return &MainConfig{
		DisplayLink: "http://localhost:8080",
		ServerHost:  "localhost:8080",
	}
}
