package subscriber

import (
	"context"
	"log/slog"
	"repo-stat/collector/internal/domain"
	proto "repo-stat/proto/subscriber"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   proto.SubscriberClient
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
		pb:   proto.NewSubscriberClient(conn),
	}, nil
}

func (c *Client) GetSubscriptionsInfo(ctx context.Context) ([]*domain.Subscription, error) {
	resp, err := c.pb.ListSubscriptions(ctx, &proto.ListSubscriptionRequest{})
	if err != nil {
		c.log.Error("failed to list subscriptions|collector", "error", err)
		return nil, err
	}

	subscriptions := make([]*domain.Subscription, 0, len(resp.Subscriptions))
	for _, sub := range resp.Subscriptions {
		subscriptions = append(subscriptions, &domain.Subscription{
			Owner: sub.Owner,
			Repo:  sub.Repo,
		})
	}

	return subscriptions, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
