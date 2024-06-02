package main

import (
	"flag"
)

func main() {
	flag.Parse()

	config := NewConfig()
	config.UpdateByOptions(options)
	config.UpdateByEnv()

	server := CreateServer(config)
	server.Start()
}
