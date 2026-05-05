-- name: CreateSubscription :exec
INSERT INTO subscriptions (owner, repo)
VALUES ($1, $2)
ON CONFLICT (owner, repo) DO NOTHING;

-- name: ListSubscriptions :many
SELECT owner, repo FROM subscriptions ORDER BY id;

-- name: DeleteAllSubscriptions :exec
TRUNCATE TABLE subscriptions RESTART IDENTITY;

-- name: GetRepoFromCache :one
SELECT * FROM repo_cache 
WHERE owner = $1 AND repo = $2;

-- name: UpsertRepoCache :exec
INSERT INTO repo_cache (owner, repo, full_name, description, stars, forks, visibility, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (owner, repo) DO UPDATE SET
    full_name         = EXCLUDED.full_name,
    description       = EXCLUDED.description,
    stars             = EXCLUDED.stars,
    forks             = EXCLUDED.forks,
    visibility        = EXCLUDED.visibility,
    created_at        = EXCLUDED.created_at
WHERE repo_cache.id IS NOT NULL; -- предотвращаем попытку обновления id