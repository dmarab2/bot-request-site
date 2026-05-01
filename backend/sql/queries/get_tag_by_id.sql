-- name: GetTagByID :one
SELECT *
FROM tags
WHERE id = $1;