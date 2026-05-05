package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type Ping struct {
	subscriber Pinger
	processor  Pinger
}

func NewPing(subscriber, processor Pinger) *Ping {
	return &Ping{
		subscriber: subscriber,
		processor:  processor,
	}
}

func (u *Ping) ProcessorPing(ctx context.Context) domain.PingStatus {
	return u.processor.Ping(ctx)
}

func (u *Ping) SubscriberPing(ctx context.Context) domain.PingStatus {
	return u.subscriber.Ping(ctx)
}
