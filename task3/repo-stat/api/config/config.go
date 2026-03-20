package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/httpserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name"`
}

type Services struct {
	Subscriber string `yaml:"subscriber"`
}

type Config struct {
	App      App
	Services Services
	HTTP     httpserver.Config
	Logger   logger.Config
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
