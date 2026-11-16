CREATE TABLE IF NOT EXISTS branches (
    id         bigserial PRIMARY KEY,
    project_id bigint NOT NULL REFERENCES projects(id),
    name       text NOT NULL,
    created_by bigint NOT NULL REFERENCES users(id),
    created_at timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT branches_name_not_blank
        CHECK (length(btrim(name)) > 0)
);

CREATE UNIQUE INDEX branches_project_name_unique
ON branches (project_id, lower(btrim(name)));
