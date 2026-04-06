-- name: CreateRequestClaim :one
INSERT INTO request_claims(request_id, claimed_at, claim_secret_hash, expires_at)
VALUES($1, NOW(), $2, NOW() + INTERVAL '30 days')
RETURNING *;