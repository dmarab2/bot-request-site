-- +goose Up
CREATE TABLE tags(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tag_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    claim_secret_hash TEXT NOT NULL
);

-- +goose Down
DROP TABLE tags;