package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"

	httpSwagger "github.com/swaggo/http-swagger"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, pingUC *usecase.Ping, repoUC *usecase.GetRepoInfo, subUC *usecase.SubscriptionUseCase) {
	mux.Handle("GET /api/ping", NewPingHandler(log, pingUC))

	mux.Handle("GET /api/repositories/info", NewRepoHandler(log, repoUC))
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	mux.Handle("POST /subscriptions", NewCreateSubscriptionHandler(log, subUC))
	mux.Handle("DELETE /subscriptions/{owner}/{repo}", NewDeleteSubscriptionHandler(log, subUC))
	mux.Handle("GET /subscriptions", NewListSubscriptionsHandler(log, subUC))

	mux.Handle("GET /subscriptions/info", NewGetSubscriptionsInfoHandler(log, repoUC))
}
