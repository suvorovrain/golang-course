package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, ping *usecase.Ping) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
}
