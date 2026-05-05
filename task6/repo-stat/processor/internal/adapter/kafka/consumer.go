package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"repo-stat/processor/internal/domain"
	"repo-stat/processor/internal/usecase"

	"github.com/segmentio/kafka-go"
)

type ResponseConsumer struct {
	reader *kafka.Reader
	repo   usecase.Repository
	log    *slog.Logger
}

func NewResponseConsumer(brokers []string, groupId string, repo usecase.Repository, log *slog.Logger) *ResponseConsumer {
	return &ResponseConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			GroupID:     groupId,
			Topic:       "repo-responses",
			StartOffset: kafka.FirstOffset,
		}),
		repo: repo,
		log:  log,
	}
}

func (c *ResponseConsumer) Start(ctx context.Context) {
	c.log.Info("starting response consumer for repo-responses")

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.log.Error("failed to fetch response message", "error", err)
			continue
		}

		var resp domain.RepoResponse
		if err := json.Unmarshal(msg.Value, &resp); err != nil {
			c.log.Error("failed to unmarshal repo response", "error", err)
			continue
		}

		c.log.Info("Received response from Kafka", "owner", resp.Owner, "repo", resp.Repo, "error", resp.Error)

		if resp.Error != "" {
			c.log.Warn("received error response from collector",
				"owner", resp.Owner,
				"repo", resp.Repo,
				"error", resp.Error)
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.log.Error("failed to commit message on error", "error", err)
			}
			continue
		}

		c.log.Info("Upserting to cache", "owner", resp.Owner, "repo", resp.Repo)

		err = c.repo.UpsertRepoCache(ctx, &domain.Repository{
			Owner:       resp.Owner,
			Repo:        resp.Repo,
			FullName:    resp.FullName,
			Description: resp.Description,
			Stars:       resp.Stars,
			Forks:       resp.Forks,
			Visibility:  resp.Visibility,
			CreatedAt:   resp.CreatedAt,
		})
		if err != nil {
			c.log.Error("failed to upsert repo cache", "owner", resp.Owner, "repo", resp.Repo, "error", err)
			continue
		} else {
			c.log.Info("successfully updated cache", "owner", resp.Owner, "repo", resp.Repo)
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.log.Error("failed to commit message", "error", err)
			}
		}
	}
}

func (c *ResponseConsumer) Close() error {
	return c.reader.Close()
}
