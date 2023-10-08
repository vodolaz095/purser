-- +goose Up
CREATE TABLE secret
(
    id         uuid NOT NULL,
    body       text,
    meta       hstore,
    created_at timestamp default now(),
    PRIMARY KEY (id)
);
CREATE INDEX secret_created_at_index ON secret (created_at DESC);

-- +goose Down
DROP TABLE secret;
