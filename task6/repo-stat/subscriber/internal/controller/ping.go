package controller

import (
	"context"
	subscriberpb "repo-stat/proto/subscriber"
)

func (h *Handler) Ping(ctx context.Context, _ *subscriberpb.PingRequest) (*subscriberpb.PingResponse, error) {
	h.log.Debug("subscriberp ping request received")

	return &subscriberpb.PingResponse{
		Reply: h.ping.Execute(ctx),
	}, nil
}
