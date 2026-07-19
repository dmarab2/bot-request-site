-- name: GetRequestsFromText :many
SELECT *
FROM requests
WHERE request_search_vector @@ websearch_to_query('english', $1)
ORDER BY ts_rank(request_search_vector, websearch_to_tsquery('english', $1)) DESC
LIMIT 10;