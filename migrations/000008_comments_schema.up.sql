CREATE TABLE IF NOT EXISTS comments (
    id           bigserial PRIMARY KEY,
    commit_id    bigint NOT NULL REFERENCES commits(id),
    author_id    bigint NOT NULL REFERENCES users(id),
    timecode_sec integer,
    body         text NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT now()
);