package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/config"
	"repo-stat/subscriber/internal/adapter/github"
	"repo-stat/subscriber/internal/adapter/repository"
	grpccontroller "repo-stat/subscriber/internal/controller"
	"repo-stat/subscriber/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

func run(ctx context.Context) error {

	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting subscriber server...")
	log.Debug("debug messages are enabled")

	dbpool, err := pgxpool.New(ctx, cfg.Database.DSN())
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("successfully connected to PostgreSQL")

	ghClient := github.NewAdapter(log)
	dbRepo := repository.NewPostgresRepository(dbpool)

	subscriptionUseCase := usecase.NewSubscriptionUseCase(dbRepo, ghClient)
	pingUseCase := usecase.NewPing()

	handler := grpccontroller.NewHandler(log, subscriptionUseCase, pingUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	subscriberpb.RegisterSubscriberServer(srv.GRPC(), handler)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

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
}
