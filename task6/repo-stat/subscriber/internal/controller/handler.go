package controller

import (
	"log/slog"
	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/internal/usecase"
)

type Handler struct {
	subscriberpb.UnimplementedSubscriberServer
	log          *slog.Logger
	subscription *usecase.SubscriptionUseCase
	ping         *usecase.Ping
}

func NewHandler(log *slog.Logger, subscriptionUseCase *usecase.SubscriptionUseCase, ping *usecase.Ping) *Handler {
	return &Handler{
		log:          log,
		subscription: subscriptionUseCase,
		ping:         ping,
	}
}
