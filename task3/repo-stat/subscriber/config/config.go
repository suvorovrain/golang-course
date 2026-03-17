package config

import (
	"task3/platform/env"
	`task3/platform/grpcserver`
	"task3/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name"`
}

type Services struct {
	API string `yaml:"subscriber"`
}

type Config struct {
	App      App
	Services Services
	HTTP     grpcserver.Config
	Logger   logger.Config
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
