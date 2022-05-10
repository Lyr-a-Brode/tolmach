package main

import (
	"github.com/caarlos0/env/v6"
	"time"
)

type Config struct {
	App App
}

type App struct {
	Port            int           `env:"HOST" envDefault:"8080"`
	EnableDebug     bool          `env:"ENABLE_DEBUG" envDefault:"false"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
}

func parseConfig() (Config, error) {
	var cfg Config

	opts := env.Options{Prefix: "AS_"}

	if err := env.Parse(&cfg, opts); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
