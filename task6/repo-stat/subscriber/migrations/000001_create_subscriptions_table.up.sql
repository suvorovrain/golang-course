-- Create subscriptions table
CREATE TABLE IF NOT EXISTS subscriptions (
    id         SERIAL PRIMARY KEY,
    owner      VARCHAR(255) NOT NULL,
    repo       VARCHAR(255) NOT NULL,

    UNIQUE (owner, repo)
);

-- Create index for fast lookups by owner + repo
CREATE INDEX ON subscriptions (owner, repo);