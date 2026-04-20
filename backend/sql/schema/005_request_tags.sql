-- +goose Up
CREATE TABLE request_tags(
    request_id BIGINT,
    tag_id BIGINT,
    PRIMARY KEY (request_id, tag_id),
    FOREIGN KEY (request_id) REFERENCES requests(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);

-- +goose Down
DROP TABLE request_tags;