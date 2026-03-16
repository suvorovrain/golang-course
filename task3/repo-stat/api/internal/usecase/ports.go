package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type SubscriberPinger interface {
	Ping(ctx context.Context) (domain.PingStatus, error)
}
