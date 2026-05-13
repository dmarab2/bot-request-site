-- name: GetTagAutocompleteList :many
SELECT name
FROM tags
WHERE name LIKE $1 || '%'

UNION

SELECT name
FROM tags
WHERE name % $1


ORDER BY post_count DESC
LIMIT 5;