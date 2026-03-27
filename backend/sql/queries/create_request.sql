-- name: CreateRequest :one
INSERT INTO requests(created_at, updated_at, request_text)
VALUES(NOW(), NOW(), $1)
RETURNING *;