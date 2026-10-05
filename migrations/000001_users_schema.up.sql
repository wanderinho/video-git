CREATE TABLE IF NOT EXISTS users (
    id            bigserial PRIMARY KEY,
    username      text UNIQUE NOT NULL,
    email         text UNIQUE NOT NULL,
    password_hash text NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now()
);