package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, ping *usecase.Ping, agu *usecase.ApiGatewayUsecase, eh *dto.ErrorHandler) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping, agu, eh))
	mux.Handle("GET /api/repositories/info", NewGetInfoRepoHandler(log, agu, eh))
}
