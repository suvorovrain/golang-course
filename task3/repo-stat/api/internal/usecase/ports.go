package usecase

import (
	"context"
	"task3/api/internal/domain"
)

type SubscriberPinger interface {
	Ping(ctx context.Context) (domain.PingStatus, error)
}