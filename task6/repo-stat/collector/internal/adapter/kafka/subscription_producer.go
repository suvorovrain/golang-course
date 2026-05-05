package kafka

import (
	"context"
	"encoding/json"
	"repo-stat/collector/internal/domain"

	"github.com/segmentio/kafka-go"
)

type SubscriptionProducer struct {
	writer *kafka.Writer
}

func NewSubscriptionProducer(brokers []string) *SubscriptionProducer {
	return &SubscriptionProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "subscription-updates",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *SubscriptionProducer) PublishSubscriptions(ctx context.Context, subs []*domain.Subscription) error {
	value, err := json.Marshal(subs)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("subscriptions"),
		Value: value,
	})
}

func (p *SubscriptionProducer) Close() error {
	return p.writer.Close()
}
