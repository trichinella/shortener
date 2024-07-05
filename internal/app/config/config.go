package config

type MainConfig struct {
	DisplayLink     string
	ServerHost      string
	FileStoragePath string
	DatabaseDSN     string
}

func newConfig() *MainConfig {
	return &MainConfig{
		DisplayLink: "http://localhost:8080",
		ServerHost:  "localhost:8080",
	}
}
