package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type Pinger interface {
	Ping(ctx context.Context) domain.PingStatus
}

type Ping struct {
	pinger Pinger
}

func NewPing(pinger Pinger) *Ping {
	return &Ping{
		pinger: pinger,
	}
}

func (u *Ping) Execute(ctx context.Context) domain.PingStatus {
	return u.pinger.Ping(ctx)
}
