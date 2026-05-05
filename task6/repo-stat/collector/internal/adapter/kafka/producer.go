// repo-stat/collector/internal/adapter/kafka/response_producer.go
package kafka

import (
	"context"
	"encoding/json"

	"repo-stat/collector/internal/domain"

	"github.com/segmentio/kafka-go"
)

type ResponseProducer struct {
	writer *kafka.Writer
}

func NewResponseProducer(brokers []string) *ResponseProducer {
	return &ResponseProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "repo-responses",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *ResponseProducer) Publish(ctx context.Context, resp domain.RepoResponse) error {
	value, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(resp.Owner + "/" + resp.Repo),
		Value: value,
	})
}

func (p *ResponseProducer) Close() error {
	return p.writer.Close()
}
