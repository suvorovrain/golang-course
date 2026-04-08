package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	processor_config "repo-stat/processor/config"
	processor_adapter "repo-stat/processor/internal/adapter"
	processor_controller "repo-stat/processor/internal/controller"
	processor_usecase "repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := processor_config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting processor server...")
	log.Debug("debug messages are enabled")

	collectorCLient, err := processor_adapter.NewCollectorClient(cfg.Services.Collector)

	if err != nil {
		return fmt.Errorf("create processor adapter : %w", err)
	}

	processorUseCase := processor_usecase.NewProcessorService(collectorCLient)

	processorHandler := processor_controller.NewProcessorService(processorUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)

	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	processorpb.RegisterProcessorServer(srv.GRPC(), processorHandler)

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
