package config

import (
	"flag"
	"os"
)

type Options struct {
	ServerHost      string
	DisplayLink     string
	FileStoragePath string
}

var BaseOptions = Options{}

func init() {
	flag.StringVar(&BaseOptions.ServerHost, "a", "localhost:8080", "Server host")
	flag.StringVar(&BaseOptions.DisplayLink, "b", "http://localhost:8080", "Link displays for user")
	flag.StringVar(&BaseOptions.FileStoragePath, "f", os.TempDir()+"/short-url-db-test.json", "File path for storage")
}
