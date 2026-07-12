-- +goose Up
ALTER TABLE requests
ADD COLUMN request_search_vector tsvector
GENERATED ALWAYS AS (
    to_tsvector('english', request_text)
) STORED;

-- +goose Down
DROP COLUMN request_search_vector;