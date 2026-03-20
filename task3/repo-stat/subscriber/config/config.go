package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name"`
}

type Services struct {
	API string `yaml:"api"`
}

type Config struct {
	App      App
	Services Services
	GRPC     grpcserver.Config
	Logger   logger.Config
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
