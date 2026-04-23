-- name: CreateRequestTagLink :one
INSERT INTO request_tags(request_id, tag_id)
VALUES ($1, $2)
RETURNING *;