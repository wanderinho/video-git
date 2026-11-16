CREATE TABLE IF NOT EXISTS projects (
    id          bigserial PRIMARY KEY,
    owner_id    bigint NOT NULL REFERENCES users(id),
    name        text NOT NULL,
    description text,
    created_at  timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT projects_name_not_blank
        CHECK (length(btrim(name)) > 0)
);

CREATE UNIQUE INDEX projects_owner_name_unique
ON projects (owner_id, lower(btrim(name)));
