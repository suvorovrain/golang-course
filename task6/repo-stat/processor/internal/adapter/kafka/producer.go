package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"repo-stat/processor/internal/domain"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	log    *slog.Logger
}

func NewProducer(brokers []string, log *slog.Logger) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "repo-requests",
			Balancer: &kafka.LeastBytes{},
		},
		log: log,
	}
}

func (p *Producer) PublishRepoRequest(ctx context.Context, owner, repo string) error {
	p.log.Info("Publishing repo request to Kafka", "owner", owner, "repo", repo)

	message := domain.RepoRequest{
		Owner: owner,
		Repo:  repo,
	}

	value, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(owner + "/" + repo),
		Value: value,
	})
	if err != nil {
		p.log.Error("Failed to publish to Kafka", "error", err)
		return err
	}

	p.log.Info("Published to Kafka successfully", "owner", owner, "repo", repo)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
