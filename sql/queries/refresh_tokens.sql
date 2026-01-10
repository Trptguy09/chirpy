-- name: AddRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at, revoked_at, created_at, updated_at)
VALUES (
    $1,
    $2,
    NOW() + INTERVAL '60 days',
    NULL,
    NOW(),
    NOW()
)
RETURNING *;

-- name: GetUserByRefreshToken :one
SELECT users.id, users.created_at, users.updated_at, users.email, users.hashed_password
FROM refresh_tokens
JOIN users
ON users.id = refresh_tokens.user_id
WHERE token = $1 AND revoked_at IS NULL AND expires_at > NOW();


-- name: RevokeRefreshToken :exec

UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW()
WHERE token = $1;
