package http

import (
	"context"
	"log/slog"
	"net/http"

	"repo-stat/api/config"
	"repo-stat/api/internal/adapter/processor"
	"repo-stat/api/internal/adapter/subscriber"
	"repo-stat/api/internal/usecase"
)

func NewHandler(ctx context.Context, log *slog.Logger, cfg config.Config) (http.Handler, error) {

	subscriberClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("cannot init subscriber adapter", "error", err)
		return nil, err
	}

	processorAdapter, err := processor.NewClient(cfg.Services.Processor, log)
	if err != nil {
		log.Error("cannot init processor adapter", "error", err)
		return nil, err
	}

	repoUseCase := usecase.NewGetRepoInfo(processorAdapter)
	pingUseCase := usecase.NewPing(subscriberClient, processorAdapter)
	subscriptionUseCase := usecase.NewSubscriptionUseCase(subscriberClient)

	mux := http.NewServeMux()
	AddRoutes(mux, log, pingUseCase, repoUseCase, subscriptionUseCase)

	log.Info("HTTP handlers initialized successfully")

	var handler http.Handler = mux
	return handler, nil
}
