CREATE TABLE IF NOT EXISTS branches (
    id         bigserial PRIMARY KEY,
    project_id bigint NOT NULL REFERENCES projects(id),
    name       text NOT NULL,
    created_by bigint NOT NULL REFERENCES users(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (project_id, name)
);