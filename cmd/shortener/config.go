package main

type MainConfig struct {
	ShortLinkHost string
	CurrentHost   string
	Protocol      string
}

func NewConfig() *MainConfig {
	return &MainConfig{
		ShortLinkHost: "localhost:8080",
		CurrentHost:   "localhost:8080",
		Protocol:      "http",
	}
}
