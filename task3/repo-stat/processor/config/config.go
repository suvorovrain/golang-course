package processor_config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type APP struct {
	Processor string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-processor"`
}
type Services struct {
	Collector string `yaml:"collector" env:"COLLECTOR_ADDRESS" env-default:"localhost:8083"`
}

type Config struct {
	App      APP               `yaml:"app"`
	Logger   logger.Config     `yaml:"logger"`
	Services Services          `yaml:"services"`
	GRPC     grpcserver.Config `yaml:"grpc"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
