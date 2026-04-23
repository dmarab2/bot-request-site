-- name: GetTagID :one
SELECT tag_id
FROM tags
WHERE tag_name LIKE $1
LIMIT 1;