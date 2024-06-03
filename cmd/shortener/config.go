package main

import (
	"os"
	"strings"
)

type MainConfig struct {
	DisplayLink string
	ServerHost  string
}

func (config *MainConfig) UpdateByOptions(o Options) {
	config.ServerHost = o.ServerHost
	config.DisplayLink = strings.Trim(o.DisplayLink, "/")
}

func (config *MainConfig) UpdateByEnv() {
	envDisplayLink := os.Getenv("BASE_URL")
	if envDisplayLink != "" {
		config.DisplayLink = envDisplayLink
	}

	envServerAddr := os.Getenv("SERVER_ADDRESS")
	if envServerAddr != "" {
		config.ServerHost = envServerAddr
	}
}

func NewConfig() *MainConfig {
	return &MainConfig{
		DisplayLink: "http://localhost:8080",
		ServerHost:  "localhost:8080",
	}
}
