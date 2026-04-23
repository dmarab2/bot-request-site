-- name: CreateTag :one
INSERT INTO tags(name, created_at, updated_at)
VALUES($1, NOW(), NOW())
RETURNING *;