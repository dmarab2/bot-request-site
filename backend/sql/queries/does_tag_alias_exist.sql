-- name: DoesTagAliasExist :one
SELECT COUNT(*)
FROM tag_aliases
WHERE name = $1;