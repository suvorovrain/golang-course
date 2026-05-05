package repository

import (
	"context"
	"fmt"
	"log/slog"

	"repo-stat/processor/internal/adapter/repository/sqlc"
	"repo-stat/processor/internal/domain"
	"repo-stat/processor/internal/usecase"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	queries *sqlc.Queries
	log     *slog.Logger
}

func NewPostgresRepository(pool *pgxpool.Pool, log *slog.Logger) usecase.Repository {
	return &postgresRepository{
		queries: sqlc.New(pool),
		log:     log,
	}
}

func (r *postgresRepository) ListSubscriptions(ctx context.Context) ([]*domain.Subscription, error) {
	r.log.Info("ListSubscriptions: querying DB")

	items, err := r.queries.ListSubscriptions(ctx)
	if err != nil {
		r.log.Error("ListSubscriptions: DB error", "error", err)
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	r.log.Info("ListSubscriptions: got from DB", "count", len(items))

	subs := make([]*domain.Subscription, len(items))
	for i, item := range items {
		subs[i] = &domain.Subscription{
			Owner: item.Owner,
			Repo:  item.Repo,
		}
	}

	return subs, nil
}

func (r *postgresRepository) ReplaceAllSubscriptions(ctx context.Context, subs []*domain.Subscription) error {
	r.log.Info("ReplaceAllSubscriptions: replacing all subscriptions", "count", len(subs))

	if err := r.queries.DeleteAllSubscriptions(ctx); err != nil {
		r.log.Error("ReplaceAllSubscriptions: failed to truncate", "error", err)
		return fmt.Errorf("failed to truncate subscriptions: %w", err)
	}

	for _, sub := range subs {
		err := r.queries.CreateSubscription(ctx, sqlc.CreateSubscriptionParams{
			Owner: sub.Owner,
			Repo:  sub.Repo,
		})
		if err != nil {
			r.log.Error("ReplaceAllSubscriptions: failed to create", "owner", sub.Owner, "repo", sub.Repo, "error", err)
			return fmt.Errorf("failed to create subscription %s/%s: %w", sub.Owner, sub.Repo, err)
		}
	}

	r.log.Info("ReplaceAllSubscriptions: replaced successfully", "count", len(subs))
	return nil
}

func (r *postgresRepository) GetRepoFromCache(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	r.log.Info("GetRepoFromCache: querying DB", "owner", owner, "repo", repo)

	row, err := r.queries.GetRepoFromCache(ctx, sqlc.GetRepoFromCacheParams{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		r.log.Error("GetRepoFromCache: DB error", "error", err)
		return nil, fmt.Errorf("failed to get repo from cache: %w", err)
	}

	r.log.Info("GetRepoFromCache: got from DB", "owner", owner, "repo", repo, "full_name", row.FullName.String)

	return &domain.Repository{
		Owner:       row.Owner,
		Repo:        row.Repo,
		FullName:    row.FullName.String,
		Description: row.Description.String,
		Stars:       row.Stars.Int32,
		Forks:       row.Forks.Int32,
		Visibility:  row.Visibility.String,
		CreatedAt:   row.CreatedAt.String,
	}, nil
}

func (r *postgresRepository) UpsertRepoCache(ctx context.Context, repo *domain.Repository) error {
	r.log.Info("UpsertRepoCache: saving to DB", "owner", repo.Owner, "repo", repo.Repo, "full_name", repo.FullName)

	err := r.queries.UpsertRepoCache(ctx, sqlc.UpsertRepoCacheParams{
		Owner:       repo.Owner,
		Repo:        repo.Repo,
		FullName:    pgtype.Text{String: repo.FullName, Valid: true},
		Description: pgtype.Text{String: repo.Description, Valid: true},
		Stars:       pgtype.Int4{Int32: repo.Stars, Valid: true},
		Forks:       pgtype.Int4{Int32: repo.Forks, Valid: true},
		Visibility:  pgtype.Text{String: repo.Visibility, Valid: true},
		CreatedAt:   pgtype.Text{String: repo.CreatedAt, Valid: true},
	})
	if err != nil {
		r.log.Error("UpsertRepoCache: DB error", "error", err)
		return fmt.Errorf("failed to upsert repo cache for %s/%s: %w", repo.Owner, repo.Repo, err)
	}

	r.log.Info("UpsertRepoCache: saved successfully", "owner", repo.Owner, "repo", repo.Repo)
	return nil
}
