-- name: GetTagAutocompleteList :many
SELECT name
FROM tags
WHERE name LIKE $1 || '%'
ORDER BY post_count DESC

UNION

SELECT name
FROM tags
WHERE name % $1


LIMIT 10;