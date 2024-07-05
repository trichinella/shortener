package config

import (
	"flag"
	"os"
	"strings"
)

type options struct {
	ServerHost      string
	DisplayLink     string
	FileStoragePath string
	DatabaseDSN     string
}

var baseOptions = options{}

func init() {
	flag.StringVar(&baseOptions.ServerHost, "a", "localhost:8080", "Server host")
	flag.StringVar(&baseOptions.DisplayLink, "b", "http://localhost:8080", "Link displays for user")
	flag.StringVar(&baseOptions.DatabaseDSN, "d", "", "DSN for database")
	flag.StringVar(&baseOptions.FileStoragePath, "f", os.TempDir()+"/short-url-db-test.json", "File path for storage")
}

func (config *MainConfig) updateByFlags(o options) {
	config.ServerHost = o.ServerHost
	config.DisplayLink = strings.Trim(o.DisplayLink, "/")
	config.FileStoragePath = o.FileStoragePath
	config.DatabaseDSN = o.DatabaseDSN
}
