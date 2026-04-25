-- name: GetTagID :one
SELECT id
FROM tags
WHERE name LIKE $1 || '%'
LIMIT 1;