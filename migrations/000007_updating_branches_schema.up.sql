ALTER TABLE branches
    ADD COLUMN head_commit_id bigint
    REFERENCES commits(id) DEFERRABLE INITIALLY DEFERRED;

