-- name: CreateTag :one
INSERT INTO tags(tag_name, created_at, updated_at)
VALUES($1, NOW(), NOW())
RETURNING *;