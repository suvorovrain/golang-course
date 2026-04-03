package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"repo-stat/platform/logger"
	"repo-stat/processor/config"
	"repo-stat/processor/internal/adapter/collector"
	grpcController "repo-stat/processor/internal/controller/grpc"
	"repo-stat/processor/internal/usecase"
	"repo-stat/proto/processor"

	"google.golang.org/grpc"
)

func run(ctx context.Context) error {
	// config
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	// logger
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting processor service...")
	log.Debug("debug messages are enabled")

	collectorClient, err := collector.NewClient(cfg.Services.Collector, log)
	if err != nil {
		log.Error("error creating collector client", "error", err)
		return err
	}
	defer func() {
		if err := collectorClient.Close(); err != nil {
			log.Error("error closing collector client", "error", err)
		}
	}()

	pingUseCase := usecase.NewPing()
	repoUseCase := usecase.NewRepo(collectorClient)

	processorServer := grpcController.NewProcessorServer(pingUseCase, repoUseCase)

	lis, err := net.Listen("tcp", cfg.GRPC.Address)
	if err != nil {
		log.Error("error creating gRPC listener", "error", err)
		return err
	}

	grpcServer := grpc.NewServer()
	processor.RegisterProcessorServiceServer(grpcServer, processorServer)

	log.Info("starting processor...")

	go func() {
		<-ctx.Done()
		log.Info("shutting down processor...")
		grpcServer.GracefulStop()
	}()

	return grpcServer.Serve(lis)
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching processor error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
	cancel()
}
