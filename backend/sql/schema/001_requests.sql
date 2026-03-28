-- +goose Up
CREATE TABLE requests(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    request_text TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE requests;