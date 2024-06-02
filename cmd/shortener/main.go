package main

import (
	"flag"
)

func main() {
	flag.Parse()

	config := NewConfig()
	config.UpdateByOptions(options)

	server := CreateServer(config)
	server.Start()
}
