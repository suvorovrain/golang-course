package subscriber

import (
	"context"
	"log/slog"
	"repo-stat/api/internal/domain"

	subscirberpb "repo-stat/proto/subscriber"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscirberpb.SubscriberClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   subscirberpb.NewSubscriberClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pb.Ping(ctx, &subscirberpb.PingRequest{})
	if err != nil {
		c.log.Error("subscriber ping failed", "error", err)
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) Close() error {
	return c.conn.Close()
}
