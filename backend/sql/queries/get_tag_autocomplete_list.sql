-- name: GetTagAutocompleteList :many
SELECT name
FROM tags
WHERE name LIKE $1 || '%'
LIMIT 3;