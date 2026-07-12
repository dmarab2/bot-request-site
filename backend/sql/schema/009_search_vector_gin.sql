-- +goose Up
CREATE INDEX request_search_idx ON requests USING GIN (request_search_vector);

-- +goose Down
DROP INDEX request_search_idx;