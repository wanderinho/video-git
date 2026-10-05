CREATE TABLE IF NOT EXISTS project_members (
    project_id bigint NOT NULL REFERENCES projects(id),
    user_id    bigint NOT NULL REFERENCES users(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (project_id, user_id)
);