package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"repo-stat/collector/internal/adapter/github"
	"repo-stat/collector/internal/domain"

	"github.com/segmentio/kafka-go"
)

type TaskConsumer struct {
	reader   *kafka.Reader
	github   *github.Adapter
	producer *ResponseProducer
	log      *slog.Logger
}

func NewTaskConsumer(brokers []string, groupId string, github *github.Adapter, producer *ResponseProducer, log *slog.Logger) *TaskConsumer {
	return &TaskConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			GroupID:     groupId,
			Topic:       "repo-requests",
			MinBytes:    1,
			MaxBytes:    10e6,
			StartOffset: kafka.FirstOffset,
		}),
		github:   github,
		producer: producer,
		log:      log,
	}
}

func (c *TaskConsumer) Start(ctx context.Context) {

	c.log.Info("started consuming from repo-requests",
		"group_id", c.reader.Config().GroupID,
		"topic", "repo-requests",
		"start_offset", "first")

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.log.Error("failed to fetch message from kafka", "error", err)
			continue
		}

		c.log.Info("received raw kafka message",
			"key", string(msg.Key),
			"value", string(msg.Value),
			"partition", msg.Partition,
			"offset", msg.Offset)

		var req domain.RepoRequest
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			c.log.Error("failed to unmarshal repo request", "error", err)
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.log.Error("failed to commit bad message", "error", err)
			}
			continue
		}

		c.log.Info("received task", "owner", req.Owner, "repo", req.Repo)

		repoInfo, err := c.github.Get(ctx, req.Owner, req.Repo)

		c.log.Info("after Get", "repoInfo_nil", repoInfo == nil, "err", err)

		c.log.Info("github response", "owner", req.Owner, "repo", req.Repo, "info", repoInfo, "error", err)

		response := domain.RepoResponse{
			Owner: req.Owner,
			Repo:  req.Repo,
		}

		if err != nil {
			response.Error = err.Error()
			c.log.Error("failed to fetch repo from github", "owner", req.Owner, "repo", req.Repo, "error", err)
		} else {
			response.FullName = repoInfo.FullName
			response.Description = repoInfo.Description
			response.Stars = repoInfo.Stars
			response.Forks = repoInfo.Forks
			response.Visibility = repoInfo.Visibility
			response.CreatedAt = repoInfo.CreatedAt
		}

		if sendErr := c.producer.Publish(ctx, response); sendErr != nil {
			c.log.Error("failed to publish response to kafka", "error", sendErr)
			continue
		} else {
			c.log.Info("successfully published response to repo-responses", "owner", req.Owner, "repo", req.Repo)
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.log.Error("failed to commit message", "error", err)
			}
		}
	}
}

func (c *TaskConsumer) Close() error {
	return c.reader.Close()
}
