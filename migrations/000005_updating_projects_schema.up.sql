ALTER TABLE projects
    ADD COLUMN default_branch_id bigint REFERENCES branches(id);
