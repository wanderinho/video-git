CREATE TABLE IF NOT EXISTS commits (
    id              bigserial PRIMARY KEY,
    branch_id       bigint NOT NULL REFERENCES branches(id),
    author_id       bigint NOT NULL REFERENCES users(id),
    parent_id       bigint REFERENCES commits(id),
    message         text NOT NULL,
    minio_key       text NOT NULL,
    duration_sec    integer,
    resolution      text,
    file_size_bytes bigint,
    preview_key     text,
    status          text NOT NULL DEFAULT 'processing'
                    CHECK (status IN ('processing', 'ready', 'failed')),
    created_at      timestamptz NOT NULL DEFAULT now()
);