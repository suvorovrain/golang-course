package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"repo-stat/processor/internal/domain"
	"repo-stat/processor/internal/usecase"

	"github.com/segmentio/kafka-go"
)

type SubscriptionConsumer struct {
	reader *kafka.Reader
	repo   usecase.Repository
	log    *slog.Logger
}

func NewSubscriptionConsumer(brokers []string, groupID string, repo usecase.Repository, log *slog.Logger) *SubscriptionConsumer {
	return &SubscriptionConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			GroupID:        groupID,
			Topic:          "subscription-updates",
			CommitInterval: 1 * time.Second,
		}),
		repo: repo,
		log:  log,
	}
}

func (c *SubscriptionConsumer) Start(ctx context.Context) {
	c.log.Info("starting subscription consumer for subscription-updates")

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				c.log.Info("context cancelled, stopping subscription consumer")
				return
			}
			c.log.Error("failed to read subscription message", "error", err)
			time.Sleep(1 * time.Second)
			continue
		}

		c.log.Info("received raw message", "key", string(msg.Key), "value", string(msg.Value))

		var subs []domain.Subscription
		if err := json.Unmarshal(msg.Value, &subs); err != nil {
			c.log.Error("failed to unmarshal subscriptions", "error", err)
			continue
		}

		c.log.Info("received subscriptions update", "count", len(subs))

		ptrSubs := make([]*domain.Subscription, len(subs))
		for i := range subs {
			ptrSubs[i] = &subs[i]
		}

		if err := c.repo.ReplaceAllSubscriptions(ctx, ptrSubs); err != nil {
			c.log.Error("failed to replace subscriptions", "error", err)
			continue
		}

		c.log.Info("successfully updated subscriptions table", "count", len(subs))
	}
}

func (c *SubscriptionConsumer) Close() error {
	return c.reader.Close()
}
