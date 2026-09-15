-- name: GetTagIdFromAlias :one
SELECT tag_id
FROM tag_aliases
WHERE id = $1
LIMIT 1;