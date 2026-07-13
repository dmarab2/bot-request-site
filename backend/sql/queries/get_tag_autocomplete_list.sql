-- name: GetTagAutocompleteList :many
SELECT *
FROM tags
WHERE name LIKE $1 || '%'

UNION

SELECT *
FROM tags
WHERE name % $1


ORDER BY post_count DESC
LIMIT 10;