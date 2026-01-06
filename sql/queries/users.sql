-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE $1 = email
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    email = $2,
    hashed_password = $3,
    updated_at = NOW() 
WHERE $1 = id
RETURNING *;

-- name: UpgradeToChirpyRed :one
UPDATE users
SET 
    is_chirpy_red = TRUE
WHERE $1 = id
RETURNING *;
