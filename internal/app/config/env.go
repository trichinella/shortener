package config

import (
	"os"
)

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
