-- +goose Up
CREATE TABLE request_claims(
    request_id BIGINT PRIMARY KEY REFERENCES requests(id) ON DELETE CASCADE,
    claimed_at TIMESTAMP NOT NULL,
    claim_secret_hash TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE request_claims;