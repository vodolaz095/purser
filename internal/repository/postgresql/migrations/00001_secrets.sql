-- +goose Up
CREATE EXTENSION IF NOT EXISTS hstore;
CREATE TABLE secret
(
    id         uuid NOT NULL default gen_random_uuid(),
    body       text,
    meta       hstore,
    created_at timestamp default now(),
    PRIMARY KEY (id)
);
CREATE INDEX secret_created_at_index ON secret (created_at DESC);

-- +goose Down
DROP TABLE secret;
DROP EXTENSION IF EXISTS hstore;
