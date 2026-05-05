package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"repo-stat/processor/config"
	"repo-stat/processor/internal/adapter/kafka"
	"repo-stat/processor/internal/adapter/repository"
	"repo-stat/processor/internal/controller"
	"repo-stat/processor/internal/usecase"
	processorProto "repo-stat/proto/processor"

	"github.com/jackc/pgx/v5/pgxpool"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting processor server...")
	log.Debug("debug messages are enabled")

	pool, err := pgxpool.New(ctx, cfg.Database.DSN())
	if err != nil {
		return fmt.Errorf("failed to connect to processor db: %w", err)
	}
	defer pool.Close()

	// Repository
	repo := repository.NewPostgresRepository(pool, log)

	// Kafka Producer for gRPC responses
	producer := kafka.NewProducer([]string{cfg.Services.Kafka}, log)
	defer producer.Close() //nolint:errcheck

	// Consumer for responses from collector
	consumer := kafka.NewResponseConsumer([]string{cfg.Services.Kafka}, "processor-response-group", repo, log)

	go consumer.Start(ctx)
	defer consumer.Close() //nolint:errcheck

	// Consumer for subscription updates from collector
	subConsumer := kafka.NewSubscriptionConsumer([]string{cfg.Services.Kafka}, "processor-subscription-group", repo, log)

	go subConsumer.Start(ctx)
	defer func() {
		if err := subConsumer.Close(); err != nil {
			log.Error("failed to close subConsumer", "error", err)
		}
	}()

	// UseCase
	getRepoUseCase := usecase.NewGetRepoUseCase(repo, producer, log)

	// Handler
	handler := controller.NewHandler(getRepoUseCase)

	// gRPC Server
	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("failed to create grpc server: %w", err)
	}

	processorProto.RegisterProcessorServer(srv.GRPC(), handler)

	log.Info("processor started", "address", cfg.GRPC.Address)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("grpc server error: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		cancel()
		os.Exit(1)
	}
	cancel()
}
