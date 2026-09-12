-- name: CreateTagAlias :one
INSERT INTO tag_aliases(name, tag_id, created_at)
VALUES($1, $2, NOW())
RETURNING *;