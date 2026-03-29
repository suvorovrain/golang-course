package http

import (
	"context"
	"log/slog"
	"net/http"
	"repo-stat/api/config"
	"repo-stat/api/internal/adapter/subscriber"
	"repo-stat/api/internal/usecase"
)

func NewHandler(ctx context.Context, log *slog.Logger, cfg config.Config) (http.Handler, error) {
	subscriberClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("cannot init subscriber adapter", "error", err)
		return nil, err
	}

	pingUseCase := usecase.NewPing(subscriberClient)

	mux := http.NewServeMux()
	AddRoutes(mux, log, pingUseCase)

	var handler http.Handler = mux
	return handler, nil
}
