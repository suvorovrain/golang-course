package config

import (
	"fmt"
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-subscriber"`
}

type Services struct {
	API string `yaml:"api" env:"API_ADDRESS" env-default:"localhost:8080"`
}

type Database struct {
	Host     string `yaml:"host"     env:"DB_HOST"     env-default:"subscriber-db"`
	Port     int    `yaml:"port"     env:"DB_PORT"     env-default:"5432"`
	User     string `yaml:"user"     env:"DB_USER"     env-default:"admin"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"password"`
	DBName   string `yaml:"dbname"   env:"DB_NAME"     env-default:"subscriber_db"`
	SSLMode  string `yaml:"sslmode"  env:"DB_SSLMODE" env-default:"disable"`
}

func (d Database) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.DBName,
		d.SSLMode,
	)
}

type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Logger   logger.Config     `yaml:"logger"`
	Database Database          `yaml:"database"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
