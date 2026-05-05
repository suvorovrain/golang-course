-- name: CreateSubscription :one
INSERT INTO subscriptions (owner, repo)
VALUES ($1, $2)
RETURNING id, owner, repo;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE owner = $1 AND repo = $2;

-- name: ListSubscriptions :many
SELECT id, owner, repo
FROM subscriptions
ORDER BY id;

-- name: GetSubscription :one
SELECT id, owner, repo
FROM subscriptions
WHERE owner = $1 AND repo = $2;

-- name: ExistsSubscription :one
SELECT EXISTS (
    SELECT 1 
    FROM subscriptions 
    WHERE owner = $1 AND repo = $2
);