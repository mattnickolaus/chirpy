-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, revoked_at, user_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NULL,
    $3
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT * FROM refresh_tokens
WHERE $1 = token
LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET
    revoked_at = $2,
    updated_at = $3
WHERE $1 = token;
