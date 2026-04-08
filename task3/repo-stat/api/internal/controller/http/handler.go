package http

import (
	"context"
	"log/slog"
	"net/http"
	"repo-stat/api/config"
	"repo-stat/api/internal/adapter/processor"
	"repo-stat/api/internal/adapter/subscriber"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

func NewHandler(ctx context.Context, log *slog.Logger, cfg config.Config) (http.Handler, error) {
	subscriberClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("cannot init subscriber adapter", "error", err)
		return nil, err
	}
	processorClient, err := processor.NewProcessorClient(cfg.Services.Processor, log)
	if err != nil {
		log.Error("cannot init processor adapter", "error", err)
		return nil, err
	}

	pingUseCase := usecase.NewPing(subscriberClient)
	apiUsecase := usecase.NewUsecaseApiGateway(processorClient)
	eh := &dto.ErrorHandler{}
	mux := http.NewServeMux()
	AddRoutes(mux, log, pingUseCase, apiUsecase, eh)

	var handler http.Handler = mux
	return handler, nil
}
