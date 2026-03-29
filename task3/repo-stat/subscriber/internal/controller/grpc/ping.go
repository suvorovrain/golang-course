package grpc

import (
	"context"
	"log/slog"
	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/internal/usecase"
)

type Server struct {
	subscriberpb.UnimplementedSubscriberServer
	log  *slog.Logger
	ping *usecase.Ping
}

func NewServer(log *slog.Logger, ping *usecase.Ping) *Server {
	return &Server{
		log:  log,
		ping: ping,
	}
}

func (s *Server) Ping(ctx context.Context, _ *subscriberpb.PingRequest) (*subscriberpb.PingResponse, error) {
	s.log.Debug("subscriberp ping request received")

	return &subscriberpb.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}
