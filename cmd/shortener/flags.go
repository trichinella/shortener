package main

import "flag"

type Options struct {
	ServerHost  string
	DisplayLink string
}

var options = Options{}

func init() {
	flag.StringVar(&options.ServerHost, "a", "localhost:8080", "Server host")
	flag.StringVar(&options.DisplayLink, "b", "http://localhost:8000", "Link displays for user")
}
