CREATE TABLE IF NOT EXISTS projects (
    id          bigserial PRIMARY KEY,
    owner_id    bigint NOT NULL REFERENCES users(id),
    name        text NOT NULL,
    description text,
    created_at  timestamptz NOT NULL DEFAULT now()
);