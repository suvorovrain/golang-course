package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/api/config"
	"repo-stat/api/internal/controller/http"
	"repo-stat/platform/httpserver"
	"repo-stat/platform/logger"

	_ "repo-stat/api/docs"
)

func run(ctx context.Context) error {
	// config
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	// logger

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting server...")
	log.Debug("debug messages are enabled")

	// handler
	handler, err := http.NewHandler(ctx, log, cfg)
	if err != nil {
		log.Error("Error creating handler", "error", err)
		return err
	}

	// server
	srv := httpserver.New(cfg.HTTP, handler)
	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run http server: %w", err)
	}
	return nil
}

// @title Repo-stat API Gateway
// @version 1.0.0
// @description API for getting info about github repo by url
// @host localhost:28080
// @BasePath /
func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching server error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
	cancel()
}
