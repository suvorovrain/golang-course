package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"repo-stat/collector/config"
	"repo-stat/collector/internal/adapter/github"
	"repo-stat/collector/internal/adapter/kafka"
	"repo-stat/collector/internal/adapter/subscriber"
	"repo-stat/collector/internal/domain"
	"repo-stat/platform/logger"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "config file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting collector server...")

	// GitHub adapter
	ghAdapter := github.NewAdapter(log)

	// Kafka Response Producer
	responseProducer := kafka.NewResponseProducer([]string{cfg.Services.Kafka})

	// Kafka Subscription Producer
	subscriptionProducer := kafka.NewSubscriptionProducer([]string{cfg.Services.Kafka})

	// Kafka Task Consumer
	taskConsumer := kafka.NewTaskConsumer(
		[]string{cfg.Services.Kafka},
		"collector-task-group",
		ghAdapter,
		responseProducer,
		log,
	)

	go startSubscriptionUpdater(ctx, cfg, ghAdapter, responseProducer, subscriptionProducer, log)

	taskConsumer.Start(ctx)

	return nil
}

func startSubscriptionUpdater(ctx context.Context, cfg config.Config, ghAdapter *github.Adapter, respProducer *kafka.ResponseProducer, subProducer *kafka.SubscriptionProducer, log *slog.Logger) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	subClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("failed to create subscriber client for updater", "error", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Debug("starting scheduled subscription update")

			subs, err := subClient.GetSubscriptionsInfo(ctx)
			if err != nil {
				log.Error("failed to get subscriptions from subscriber", "error", err)
				continue
			}

			if err := subProducer.PublishSubscriptions(ctx, subs); err != nil {
				log.Error("failed to publish subscriptions update", "error", err)
			} else {
				log.Info("published subscriptions to kafka", "count", len(subs))
			}

			for _, sub := range subs {
				repoInfo, err := ghAdapter.Get(ctx, sub.Owner, sub.Repo)
				if err != nil {
					log.Error("failed to get repo from github", "owner", sub.Owner, "repo", sub.Repo, "error", err)
					continue
				}

				response := domain.RepoResponse{
					Owner:       sub.Owner,
					Repo:        sub.Repo,
					FullName:    repoInfo.FullName,
					Description: repoInfo.Description,
					Stars:       repoInfo.Stars,
					Forks:       repoInfo.Forks,
					Visibility:  repoInfo.Visibility,
					CreatedAt:   repoInfo.CreatedAt,
				}

				if err := respProducer.Publish(ctx, response); err != nil {
					log.Error("failed to publish repo response", "owner", sub.Owner, "repo", sub.Repo, "error", err)
				}
			}

			log.Info("scheduled update completed", "subscriptions_count", len(subs))
		}
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "collector error: %v\n", err)
		cancel()
		os.Exit(1)
	}
	cancel()
}
