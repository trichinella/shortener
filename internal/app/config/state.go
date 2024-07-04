package config

import (
	"flag"
	"sync"
)

var once sync.Once
var cfg *MainConfig

func State() *MainConfig {
	once.Do(func() {
		flag.Parse()

		cfg = newConfig()
		cfg.updateByOptions(baseOptions)
		cfg.updateByEnv()
	})

	return cfg
}
