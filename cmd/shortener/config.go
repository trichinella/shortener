package main

import "strings"

type MainConfig struct {
	DisplayLink string
	ServerHost  string
}

func (config *MainConfig) UpdateByOptions(o Options) {
	config.ServerHost = o.ServerHost
	config.DisplayLink = strings.Trim(o.DisplayLink, "/")
}

func NewConfig() *MainConfig {
	return &MainConfig{
		DisplayLink: "http://localhost:8080",
		ServerHost:  "localhost:8080",
	}
}
