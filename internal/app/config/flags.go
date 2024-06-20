package config

import "flag"

type Options struct {
	ServerHost  string
	DisplayLink string
}

var BaseOptions = Options{}

func init() {
	flag.StringVar(&BaseOptions.ServerHost, "a", "localhost:8080", "Server host")
	flag.StringVar(&BaseOptions.DisplayLink, "b", "http://localhost:8080", "Link displays for user")
}
