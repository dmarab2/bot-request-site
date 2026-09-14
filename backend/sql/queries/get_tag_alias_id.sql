-- name: GetTagAliasID :one
SELECT id
FROM tag_aliases
WHERE name LIKE $1 || '%'
LIMIT 1;